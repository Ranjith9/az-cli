package environment

import (
        "az-cli/azure/interface/computeinterface"
        "az-cli/azure/interface/networkinterface"
        "fmt"
)

type DeleteIN struct {
        VmName string
}

func (v DeleteIN) Delete() {

        m := azurecompute.VMIn{ResourceGroup: "CLI-group", VmName: v.VmName}
        vm, _ := m.GetVM()

        osDisk := *vm.VirtualMachineProperties.StorageProfile.OsDisk.Name
        fmt.Println(osDisk)



//************************@@@VM DELETE@@@***********************
        vd, _ := m.DeleteVM()
        if vd.Response.StatusCode == 200 { fmt.Println("Deleted VM "+m.VmName+" successfully")}

//***********************@@@NIC DELETE@@@*****cmd/create.go******************
        n := azurenetwork.NicIn{ResourceGroup: "CLI-group", NicName: v.VmName + "-nic"}
        nd, _ := n.DeleteNIC()
        if nd.Response.StatusCode == 200 { fmt.Println("Deleted NIC "+n.NicName+" successfully")}

//***********************@@@NSG DELETE@@@***********************
        s := azurenetwork.NsgIn{ResourceGroup: "CLI-group", NsgName: v.VmName + "-nsg"}
        sd, _ := s.DeleteNetworkSecurityGroup()
        if sd.Response.StatusCode == 200 { fmt.Println("Deleted NSG "+s.NsgName+" successfully")}

//************************@@@IP DELETE@@@***********************
        i := azurenetwork.IpIn{ResourceGroup: "CLI-group", IpName: v.VmName + "-ip"}
        id, _ := i.DeletePublicIP()
        if id.Response.StatusCode == 200 { fmt.Println("Deleted IP "+i.IpName+" successfully")}

//***********************@@@DISK DELETE@@***********************
        d := azurecompute.DisksIn{ResourceGroup: "CLI-group", DiskName: osDisk}
        dd,_ := d.DeleteDisk()
        if dd.Response.StatusCode == 200 { fmt.Println("Deleted DISK "+d.DiskName+" successfully")}
}
