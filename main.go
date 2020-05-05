package main

import (
	"fmt"
	"main/model"
	"main/services"
	"main/services/api"
	"strings"
)

// 在切片中查找元素
func findInSlice(length int, f func(i int) bool) int {

	for i := 0; i < length; i++ {
		if f(i) {
			return i
		}
	}

	return -1
}

func deleteHostInCloudFlare(config *model.Config, recordlist []model.DNSRecord, hosts []string) {

	for _, record := range recordlist {
		fmt.Printf("当前域名为: %s , ip:  %s , type: %s \n", record.Name, record.Content, record.Type)

		findIndex := findInSlice(len(hosts), func(i int) bool {
			fmt.Println(hosts[i], record.Name, strings.EqualFold(hosts[i], record.Name))
			if strings.EqualFold(hosts[i], record.Name) {
				return true
			}

			return false
		})

		// 符合条件的清除
		if findIndex != -1 {
			switch record.Type {
			case "A":
				fallthrough
			case "AAAA":
				id, err := api.DeleteDNSRecord(config, record.ID)
				if err != nil {
					fmt.Println(err)
					continue
				}
				fmt.Printf("删除成功id: %s", id)
			default:
			}
		}

	}
}

func createDDNS(config *model.Config, ips services.IPs, zoneID string, ipv6NeighborList []services.Ipv6NeighborInfo) {
	for _, dns := range config.Ddns {
		host := dns.Host
		mac := dns.Mac

		if dns.IsCurrent { // 当前路由器
			// 当前路由器的ddns配置
			if len(ips.IPv4) != 0 {
				// 有ipv4时创建dns
				result, err := api.CreateDNSRecord(config, zoneID, "A", host, ips.IPv4[0], true)

				if err != nil {
					fmt.Printf("创建失败 %s", err)
					continue
				}

				fmt.Printf("创建成功, %#v", result)
			} else if len(ips.IPv6) != 0 {
				fmt.Printf("length for ipv6 :%v , %v \n\n", len(ips.IPv6), ips.IPv6)
				// 有ipv6时创建dns
				result, err := api.CreateDNSRecord(config, zoneID, "AAAA", host, ips.IPv6[0], true)

				if err != nil {
					fmt.Printf("创建失败 %s", err)
					continue
				}

				fmt.Printf("创建成功, %#v", result)
			}

		} else if mac != "" { // 路由器下面的子主机并且mac地址存在

			findMacIndex := findInSlice(len(ipv6NeighborList), func(i int) bool {
				if strings.EqualFold(ipv6NeighborList[i].Mac, mac) {
					return true
				}
				return false
			})

			if findMacIndex != -1 {

				// 创建dns
				result, err := api.CreateDNSRecord(config, zoneID, "AAAA", dns.Host, ipv6NeighborList[findMacIndex].Addr, true)

				if err != nil {
					fmt.Printf("创建失败 %s", err)
					continue
				}

				fmt.Printf("创建成功, %#v", result)

			}

		}
	}
}

func task() {

	config, err := services.LoadConfig("./config.json")

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%v", config)

	zoneList, err := api.GetZoneRecord(config)

	if err != nil {
		fmt.Println(err)
		return
	}

	if len(zoneList) == 0 {
		fmt.Println("当前用户没有zone列表，请检查cloudflare")
		return
	}

	findIndex := findInSlice(len(zoneList), func(i int) bool {
		if strings.EqualFold(zoneList[i].Name, config.DomainName) {
			return true
		}

		return false
	})

	if findIndex == -1 {
		fmt.Println("在cloudflare未找到当前域名，任务结束")
		return
	}

	zoneID := zoneList[findIndex].ID

	recordlist, err := api.GetDNSRecord(config, zoneID, "AAAA")

	if err != nil {
		fmt.Println(err)

		return
	}

	hosts := config.GetConfigAllHost()

	// 先删除包含在内的域名
	deleteHostInCloudFlare(config, recordlist, hosts)

	ips, err := services.GetCurrentIPs()

	if err != nil {
		fmt.Println(err)
	}

	ipv6NeighborList, err := services.GetIpv6NeighborList()

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(ipv6NeighborList)

	createDDNS(config, ips, zoneID, ipv6NeighborList)

}

func main() {
	task()
}
