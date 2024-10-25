package consul

import (
	"fmt"
	"strings"

	"github.com/hashicorp/consul/api"
	"github.com/hedzr/cmdr"
	"github.com/hedzr/consul-tags/objects/consul"
)

func ModifyServiceTagByID(id string, newTag string) {
	registrar := getRegistrar(ConsulAddr, Scheme, Token)
	modifyServiceTagByID(registrar, id, newTag)
	return
}

func ListTagsByName(name string) (err error) {
	registrar := getRegistrar(ConsulAddr, Scheme, Token)
	listServiceTagsByName(registrar, name)
	return
}

func FilterTag(name string, tag string) (err error) {
	registrar := getRegistrar(ConsulAddr, Scheme, Token)
	filterTag(registrar, name, tag)
	return
}

func listServiceTagsByName(registrar *consul.Registrar, serviceName string) {
	cmdr.Logger.Debugf("List the tags of service '%s' at '%s'...", serviceName, ConsulAddr)

	catalogServices, err := consul.QueryService(serviceName, registrar.FirstClient.Catalog())
	if err != nil {
		fmt.Printf("QueryService err: %s", err.Error())
		return
	}

	if _, ok := catalogServices[0].TaggedAddresses["wan"]; !ok {
		fmt.Printf("Bad: %v\n", catalogServices[0])
		return
	}

	for _, catalogService := range catalogServices {
		fmt.Printf("%s:\n", catalogService.ServiceID)
		fmt.Printf("\tname: %s\n", catalogService.ServiceName)
		fmt.Printf("\tnode: %s\n", catalogService.Node)
		fmt.Printf("\taddr: %s, tagged: %v\n", catalogService.Address, catalogService.TaggedAddresses)
		fmt.Printf("\tendp: %s:%d\n", catalogService.ServiceAddress, catalogService.ServicePort)
		fmt.Printf("\ttags: %v\n", strings.Join(catalogService.ServiceTags, ","))
		fmt.Printf("\tmeta: %v\n", catalogService.NodeMeta)
		fmt.Printf("\tenableTagOveerride: %v\n", catalogService.ServiceEnableTagOverride)
	}
}

func filterTag(registrar *consul.Registrar, serviceName, tag string) {
	cmdr.Logger.Debugf("List the tags of service '%s' at '%s'...", serviceName, ConsulAddr)

	catalogServices, err := consul.QueryService(serviceName, registrar.FirstClient.Catalog())
	if err != nil {
		fmt.Printf("QueryService err: %s", err.Error())
		return
	}

	if _, ok := catalogServices[0].TaggedAddresses["wan"]; !ok {
		fmt.Printf("Bad: %v\n", catalogServices[0])
		return
	}

	for _, catalogService := range catalogServices {
		if catalogService.ServiceTags[0] != tag {
			continue
		}
		fmt.Printf("id: %s\n", catalogService.ServiceID)
		fmt.Printf("\tname: %s\n", catalogService.ServiceName)
		fmt.Printf("\tnode: %s\n", catalogService.Node)
		fmt.Printf("\taddr: %s, tagged: %v\n", catalogService.Address, catalogService.TaggedAddresses)
		fmt.Printf("\tendp: %s:%d\n", catalogService.ServiceAddress, catalogService.ServicePort)
		fmt.Printf("\ttags: %v\n", strings.Join(catalogService.ServiceTags, ","))
		fmt.Printf("\tmeta: %v\n", catalogService.NodeMeta)
		fmt.Printf("\tenableTagOveerride: %v\n", catalogService.ServiceEnableTagOverride)
	}

	return
}

func modifyServiceTagByID(registrar *consul.Registrar, id string, newTag string) (err error) {
	cmdr.Logger.Debugf("Modifying the tags of service by id '%s'...", id)

	as0, err := consul.QueryServiceByID(id, registrar.FirstClient)
	if err != nil {
		cmdr.Logger.Errorf("Error: %v", err)
		return
	}

	s, err := consul.AgentServiceToCatalogService(as0, registrar.FirstClient)
	if err != nil {
		cmdr.Logger.Errorf("Error: %v", err)
		return
	}

	// 服务 s 所在的 Node
	cn := consul.NodeToAgent(registrar, s.Node)

	client := GetClient(ConsulAddr, Scheme, Token)

	err = client.Agent().ServiceRegister(&api.AgentServiceRegistration{
		ID:                as0.ID,
		Name:              as0.Service,
		Tags:              []string{newTag},
		Port:              as0.Port,
		Address:           cn.Services[id].Address,
		EnableTagOverride: as0.EnableTagOverride,
	})
	if err != nil {
		cmdr.Logger.Errorf("Error: %v", err)
		return
	}

	// 重新载入 s 的等价物才能得到新的 tags 集合，s.ServiceTags 并不会自动更新为新集合
	sNew, err := consul.QueryServiceByID(as0.ID, client)
	if err != nil {
		cmdr.Logger.Errorf("Error: %v", err)
		return
	}

	fmt.Printf("%s: %s\n", s.ServiceID, strings.Join(sNew.Tags, ","))
	return
}

func getRegistrar(addr, scheme, token string) *consul.Registrar {
	return &consul.Registrar{
		Base: consul.Base{
			FirstClient: consul.MakeClientWithConfig(func(clientConfig *api.Config) {
				clientConfig.Address = addr
				clientConfig.Scheme = scheme
				clientConfig.Token = token
			}),
		},
		Clients:       nil,
		CurrentClient: nil,
	}
}

func GetClient(addr, scheme, token string) *api.Client {
	return consul.MakeClientWithConfig(func(clientConfig *api.Config) {
		clientConfig.Address = addr
		clientConfig.Scheme = scheme
		clientConfig.Token = token
	})
}
