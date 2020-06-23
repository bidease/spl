package main

import (
	"fmt"
	"log"
	"os"

	"github.com/bidease/spl"
	"github.com/bidease/spl/cloud"
	"github.com/bidease/spl/dedicated"
	"github.com/olekukonko/tablewriter"
)

func printSSHKeys() {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"name", "created at", "updated at", "fingerprint"})

	sshkeys, err := spl.GetSSHKeys()
	if err != nil {
		log.Fatalln(err)
	}

	for _, item := range *sshkeys {
		table.Append([]string{
			item.Name,
			item.CreatedAt,
			item.UpdatedAt,
			item.Fingerprint,
		})
	}

	table.Render()
}

func printLocations() {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"id", "name", "l2 segments", "private racks", "load balancers"})

	locations, err := spl.GetLocations()
	if err != nil {
		log.Fatalln(err)
	}

	for _, item := range *locations {
		table.Append([]string{
			fmt.Sprint(item.ID),
			item.Name,
			fmt.Sprint(item.L2SegmentsEnabled),
			fmt.Sprint(item.PrivateRacksEnabled),
			fmt.Sprint(item.LoadBalancersEnabled),
		})
	}

	table.Render()
}

func printDedicatedServers() {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"id", "title", "location id", "status", "private ipv4 address", "public ipv4 address"})

	dHosts, err := dedicated.GetServers()
	if err != nil {
		log.Fatalln(err)
	}

	for _, item := range *dHosts {
		table.Append([]string{
			item.ID,
			item.Title,
			fmt.Sprint(item.LocationID),
			item.Status,
			item.PrivateIPv4Address,
			item.PublicIPv4Address,
		})
	}

	table.Render()
}

func printCloudInstances() {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"id", "name", "status", "flavor id", "private ipv4 address", "public ipv4 address"})

	cInstance, err := cloud.GetInstances()
	if err != nil {
		log.Fatalln(err)
	}

	for _, item := range *cInstance {
		table.Append([]string{
			item.ID,
			item.Name,
			item.Status,
			item.FlavorID,
			item.PrivateIPv4Address,
			item.PublicIPv4Address,
		})
	}

	table.Render()
}

func getDedicatedServersDescribe() {
	server, err := dedicated.GetServer(conf.Hostid)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println()
	fmt.Println("ID:                   ", server.ID)
	fmt.Println("Title:                ", server.Title)
	fmt.Println("Created at:           ", server.CreatedAt)
	fmt.Println("Updated at:           ", server.UpdatedAt)
	fmt.Println("Lease start at:       ", server.LeaseStartAt)
	fmt.Println("Scheduled release at: ", server.ScheduledReleaseAt)
	fmt.Println("Type:                 ", server.Type)
	fmt.Println("Location ID:          ", server.LocationID)
	fmt.Println("Status:               ", server.Status)
	fmt.Println("Configuration:        ", server.Configuration)
	fmt.Println("Private IPv4 address: ", server.PrivateIPv4Address)
	fmt.Println("Public IPv4 address:  ", server.PublicIPv4Address)
	fmt.Println()
}

func getCloudInstanceDescribe() {
	instance, err := cloud.GetInstance(conf.Hostid)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println()
	fmt.Println("ID:                   ", instance.ID)
	fmt.Println("Name:                 ", instance.Name)
	fmt.Println("Openstack UUID:       ", instance.OpenstackUUID)
	fmt.Println("Created at:           ", instance.CreatedAt)
	fmt.Println("Updated at:           ", instance.UpdatedAt)
	fmt.Println("Status:               ", instance.Status)
	fmt.Println("Flavor ID             ", instance.FlavorID)
	fmt.Println("Image ID:             ", instance.ImageID)
	fmt.Println("Private IPv4 address: ", instance.PrivateIPv4Address)
	fmt.Println("Public IPv4 address:  ", instance.PublicIPv4Address)
	fmt.Println("Public IPv6 address:  ", instance.PublicIPv6Address)
	fmt.Println()
}
