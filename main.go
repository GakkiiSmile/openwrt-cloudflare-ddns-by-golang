package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"main/model"
	"main/services"
	"main/services/api"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var currentPwd, _ = filepath.Abs(filepath.Dir(os.Args[0]))

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
	tmp := &tmpJSON{}

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

				tmp.Ddns = append(tmp.Ddns, tmpDDNS{
					Host: host,
					IP:   ips.IPv4[0],
				})

				fmt.Printf("创建成功, %#v", result)
			} else if len(ips.IPv6) != 0 {
				fmt.Printf("length for ipv6 :%v , %v \n\n", len(ips.IPv6), ips.IPv6)
				// 有ipv6时创建dns
				result, err := api.CreateDNSRecord(config, zoneID, "AAAA", host, ips.IPv6[0], true)

				if err != nil {
					fmt.Printf("创建失败 %s", err)
					continue
				}

				tmp.Ddns = append(tmp.Ddns, tmpDDNS{
					Host: host,
					IP:   ips.IPv6[0],
				})

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

				tmp.Ddns = append(tmp.Ddns, tmpDDNS{
					Host: dns.Host,
					IP:   ipv6NeighborList[findMacIndex].Addr,
				})

				fmt.Printf("创建成功, %#v", result)

			}

		}
	}

	saveTmpFile(tmp)
}

type tmpDDNS struct {
	Host string `json:"host"`
	IP   string `json:"ip"`
}

// 临时文件结构体
type tmpJSON struct {
	Ddns []tmpDDNS `json:"ddns"`
}

func readTempFile() (jsonData *tmpJSON, err error) {
	tmpFilePath := path.Join(currentPwd, "./cloudflare.ddns.tmp.json")
	fileObj, err := os.OpenFile(tmpFilePath, os.O_CREATE|os.O_RDWR, 0x644)

	if err != nil {
		fmt.Println("文件读取失败", err)
		return jsonData, err
	}
	defer fileObj.Close()

	file, err := ioutil.ReadAll(fileObj)

	if err != nil {
		fmt.Println("文件读取失败", err)
		return
	}

	jsonData = &tmpJSON{}

	if string(file) != "" {
		err = json.Unmarshal(file, jsonData)

		if err != nil {
			fmt.Println("临时文件json解析失败", err)
			return jsonData, err
		}
	}

	return jsonData, err
}

// 判断是否有相同的记录
func hasSameRecordInTempFile(oldData *tmpJSON, recordlist []model.DNSRecord) bool {
	if len(oldData.Ddns) > 0 {

		findIndex := findInSlice(len(oldData.Ddns), func(i int) bool {

			findIndex := findInSlice(len(recordlist), func(j int) bool {
				if recordlist[j].Name == oldData.Ddns[i].Host && recordlist[j].Content == oldData.Ddns[i].IP {
					return true
				}

				return false
			})

			if findIndex != -1 {
				return true
			}

			return false

		})

		// 如果能找到一个一样的，则数据没有变化，直接退出
		if findIndex != -1 {
			return true
		}
	}

	return false
}

func saveTmpFile(tmp *tmpJSON) {
	data, err := json.Marshal(tmp)

	if err != nil {
		return
	}

	tmpFilePath := path.Join(currentPwd, "./cloudflare.ddns.tmp.json")
	err = ioutil.WriteFile(tmpFilePath, data, 0x644)

	if err != nil {
		return
	}

	fmt.Println("历史ddns记录已保存", err)
}

func task() {
	configFilePath := path.Join(currentPwd, "./config.json")

	config, err := services.LoadConfig(configFilePath)

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

	// 读取本地ip
	ips, err := services.GetCurrentIPs()

	if err != nil {
		fmt.Println(err)
	}

	// 读取邻居ip
	ipv6NeighborList, err := services.GetIpv6NeighborList()

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(ipv6NeighborList)

	// 先判断当前dns是否已经发生变化，发生变化再继续往下执行替换
	oldData, err := readTempFile()

	if err != nil {
		fmt.Println("读取临时文件出错")
		return
	}

	// 如果临时文件，没有数据记录，则直接往下执行，有数据则对比 clouflare 的数据
	hasSame := hasSameRecordInTempFile(oldData, recordlist)
	if hasSame {
		return
	}

	// 先删除包含在内的域名
	deleteHostInCloudFlare(config, recordlist, hosts)

	createDDNS(config, ips, zoneID, ipv6NeighborList)

}

func main() {
	task()
}
