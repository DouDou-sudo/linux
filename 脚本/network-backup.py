[2023-05-28 11:25:04][PID:488513] job_start
#!/var/lib/sdsom/venv/bin/python
# -*- coding: utf8 -*-
import commands
import json
import os
import sys
import re
import configobj

sign = 0 
ifcfg_eth_prefix="/etc/sysconfig/network-scripts/ifcfg-"
ifcfg_path="/etc/sysconfig/network-scripts/"
unuse_ifcfg_file_eth_list = []
unuse_ifcfg_file_bond_list = []

def ten_to_hex(num):
    result = str(hex(num)).replace("0x","")
    if len(result) == 1:
        result = "0"+result
    return result


def get_ip():
    ret,output = commands.getstatusoutput("sudo ip a|grep glo|grep -v second|grep -E \"bond[0-2](\.[0-9]*)?$\"|grep -oE \"([0-9]{1,3}\.){3}[0-9]{1,3}/[0-9]{1,}\"")
    ip_list = output.split("\n")
    ip_eth_info = {}
    for ip in ip_list:
        ret,output = commands.getstatusoutput("sudo ip a|grep glo|grep {}|grep -oE \"bond[0-2](\.[0-9]*)?$\"".format(ip))
        ip_eth_info[ip] = output 
    return ip_eth_info

def get_ipv6(ip):
    ret,output = commands.getstatusoutput("sudo ip a|grep glo|grep inet6|awk '{print $2}'|awk -F \"::\" '{print $1}'|sort|uniq")
    ipv6 = ""
    ipv6_prefix_cut = []
    ipv6_prefix = ""
    if len(output.split("\n")) == 1:
        ipv6_prefix_cut = output.split("\n")[0].split(":")
        for i in range(len(ipv6_prefix_cut)):
            if len(ipv6_prefix_cut[i]) == 3:
                ipv6_prefix_cut[i] = "0"+ipv6_prefix_cut[i]
        len_tap = len(ipv6_prefix_cut)-1
        for i in range(len_tap):
            ipv6_prefix = ipv6_prefix+ipv6_prefix_cut[i]+":"
        ipv6_prefix = ipv6_prefix+ipv6_prefix_cut[len_tap]
        ip_cut = re.split("[./]",ip)
        if int(ip_cut[4]) == 24:
            ipv6_mask = "/120"
        elif int(ip_cut[4]) == 26:
            ipv6_mask = "/122"
        elif int(ip_cut[4]) == 28:
            ipv6_mask = "/124"
        elif int(ip_cut[4]) == 30:
            ipv6_mask = "/126"
        elif int(ip_cut[4]) == 32:
            ipv6_mask = "/128"
        else:
            return ipv6
        
        ipv6 = ipv6_prefix+"::"+ten_to_hex(int(ip_cut[0]))+ten_to_hex(int(ip_cut[1]))+":"+ten_to_hex(int(ip_cut[2]))+ten_to_hex(int(ip_cut[3]))+ipv6_mask
    else:
        return ipv6  
    return ipv6.upper()

def fix_ifcfg_file_ipv6(eth,ip,ipv6):
    global ifcfg_eth_prefix
    global sign
    ifcfg_eth_path = ifcfg_eth_prefix+eth
    if os.path.exists(ifcfg_eth_path) and os.path.isfile(ifcfg_eth_path):
        try:
            conf = configobj.ConfigObj(ifcfg_eth_path)
        except:
            sign = 1
            print "[ERROR]ifcfg_has_dumplicated_key_or_format_err:",ifcfg_eth_path
            return
        try:
            ipv6_config = conf["IPV6ADDR"].upper()
            if ipv6_config == ipv6:
                print "[INFO]ipv6_addr_is_ok:",ifcfg_eth_path,ipv6_config
            else:
                print "[INFO]ipv6_addr_is_notok:",ifcfg_eth_path,"config:",ipv6_config,"ipv4:",ip,"ipv6_should_be:",ipv6
                cmd1 = "sudo sed -i \"/IPV6ADDR/d\" {}".format(ifcfg_eth_path)
                cmd2 = "sudo sed -i \"/IPV6_ADDR_GEN_MODE/aIPV6ADDR={}\" {}".format(ipv6,ifcfg_eth_path)
                print "[fix]now_fix:change",cmd1,cmd2
                ret1,output1 = commands.getstatusoutput(cmd1)
                ret2,output2 = commands.getstatusoutput(cmd2)
        except:
            print "[INFO]ifcfg_has_no_ipv6addr:",ifcfg_eth_path
        return

def fix_ifcfg_special_word(ifcfg_file):
    global ifcfg_path
    ifcfg_file_path = ifcfg_path+ifcfg_file
    print "[INFO]fix_route_file_excel_special_word"
    os.system("sudo sed -i 's/\xc2\xa0/ /g' {}".format(ifcfg_file_path))
    os.system("sudo sed -i 's/\t/ /g' {}".format(ifcfg_file_path))
    os.system("sudo sed -i 's/ \+/ /g' {}".format(ifcfg_file_path))
    os.system("sudo sed -i 's/^ \*#/#/g' {}".format(ifcfg_file_path))
    return

def reduce_ifcfg_file_sameline(ifcfg_file):
    global ifcfg_path
    ifcfg_file_path = ifcfg_path+ifcfg_file
    tmp_file_path = "/tmp/ifcfg_new"
    line_list = []
    f = open(ifcfg_file_path,'r')
    f_new = open(tmp_file_path,'w')
    for line in f.readlines():
        if line not in line_list:
            line_list.append(line)
            f_new.write(line)
        else:
            continue
    f_new.close()
    os.chmod(tmp_file_path,0766)
    print "[INFO]now_reduce_same_line:",ifcfg_file_path
    os.system("sudo cat {} >{}".format(tmp_file_path,ifcfg_file_path))
    return
    
def get_and_bak_ifcfg_file():
    global ifcfg_path
    global sign
    ret,output = commands.getstatusoutput("sudo ls  %s|grep \"ifcfg-\"" %ifcfg_path)
    ifcfg_file_list = output.split("\n")
    
    for file in ifcfg_file_list:
        file_path =  ifcfg_path + file 
        if not os.path.isfile(file_path):
            ifcfg_file_list.remove(file)
    if os.path.exists(ifcfg_path+"ifcfgfixbak_2/"):
        print "[INFO]{}_already_exists_noneed_bak".format(ifcfg_path+"ifcfgfixbak_2/")
        return ifcfg_file_list
    else:
        os.system("sudo mkdir {}".format(ifcfg_path+"ifcfgfixbak_2/"))
    for file in ifcfg_file_list:
        cmd = "sudo cp -ar {} {}".format(ifcfg_path+file,ifcfg_path+"ifcfgfixbak_2/")
        print "[INFO]now_bak:",cmd
        ret,output = commands.getstatusoutput(cmd)
        if ret != 0 :
            sign = 1
            print "[ERROR]bak_fail:",cmd
    return ifcfg_file_list
 
def get_bond_list():
    ret,output = commands.getstatusoutput("sudo ip a|grep -oE \"bond[0-2]\"|sort|uniq")
    return output.split("\n")

def get_bond_eth_list():
    ret,output = commands.getstatusoutput("sudo ip a|grep -oE \"bond[0-2](\.[0-9]*)?(:21)?(:22)?$\"")
    bond_id = get_bond_list()
    bond_eth_list = output.split("\n")+bond_id
    return bond_eth_list

def get_bond_slave(bond_id):
    slave_for_bond = {}
    ret,output = commands.getstatusoutput("sudo cat /proc/net/bonding/%s|grep -i interface|awk '{print $NF}'" %(bond_id))
    for i in output.split("\n"):
        slave_for_bond[i] = bond_id
    return slave_for_bond

def check_ifexists_ifcfg_bond(bond_eth_list):
    global sign
    global ifcfg_eth_prefix
    global ifcfg_path
    ret,output = commands.getstatusoutput("sudo ls  %s|grep \"ifcfg-\"|awk -F \"ifcfg-\" '{print $NF}'" %ifcfg_path)
    ifcfg_file_list = output.split("\n")
    print "[INFO]network-scripsts目录下存在ifcfg文件:",json.dumps(ifcfg_file_list).replace(', ','+')
    for eth_id in bond_eth_list:
        ifcfg_file_path = ifcfg_eth_prefix+eth_id
        if not os.path.exists(ifcfg_file_path):
            sign = 1
            print "[ERROR][ifcfg_bond_file_check]",ifcfg_file_path," not_exists!"
    return

def check_unuse_ifcfg_bond(bond_eth_list):
    global ifcfg_eth_prefix
    global sign
    global unuse_ifcfg_file_bond_list
    #ret,output = commands.getstatusoutput("sudo ls  %sbond*|grep -E \"bond[0-2](\.[0-9]*)?(:21)?(:22)?$\"|awk -F \"ifcfg-\" '{print $NF}'" %ifcfg_eth_prefix)
    ret,output = commands.getstatusoutput("sudo ls  %s|grep -E \"ifcfg-bond[0-2](\.[0-9]*)?(:21)?(:22)?$\"|awk -F \"ifcfg-\" '{print $NF}'" %ifcfg_path)
    ifcfg_bond_list = output.split("\n")
    #print ifcfg_bond_list
    for i in ifcfg_bond_list:
        if i not in bond_eth_list and os.path.isfile(ifcfg_eth_prefix+i):
            unuse_ifcfg_file_bond_list.append(i)
            #print "unuse_ifcfg_bond_file:",ifcfg_eth_prefix+i
            ##sign = 1
    return

def check_unuse_ifcfg_interface(slave_for_bond):
    global ifcfg_eth_prefix
    global sign
    unuse_ifcfg_file_eth_list = []
    #ret,output = commands.getstatusoutput("sudo ls  %s*|grep -vE \"bond[0-2](\.[0-9]*)?(:21)?(:22)?$\"|awk -F \"ifcfg-\" '{print $NF}'" %ifcfg_eth_prefix)
    ret,output = commands.getstatusoutput("sudo ls  %s|grep \"ifcfg-\"|grep -vE \"bond[0-2](\.[0-9]*)?(:21)?(:22)?$\"|awk -F \"ifcfg-\" '{print $NF}'" %ifcfg_path)
    ifcfg_interface_list = output.split("\n")
    for i in ifcfg_interface_list:
        if not os.path.isfile(ifcfg_eth_prefix+i):
            ifcfg_interface_list.remove(i)
    print "[INFO]network-scripsts目录下存在物理网卡文件:",json.dumps(ifcfg_interface_list).replace(', ','+')
    for i in ifcfg_interface_list:
        if i != "lo" and i not in slave_for_bond.keys() and os.path.isfile(ifcfg_eth_prefix+i):
            unuse_ifcfg_file_eth_list.append(i)
    return unuse_ifcfg_file_eth_list

def check_lost_ifcfg_interface(slave_for_bond):
    global ifcfg_eth_prefix
    global sign
    lose_ifcfg_eth_list = []
    for i in slave_for_bond.keys():
        if not os.path.exists(ifcfg_eth_prefix+i):
            lose_ifcfg_eth_list.append(i)
    return lose_ifcfg_eth_list

def from_unuse_ifcfg_file_find_eth_config(eth,bond):
    global unuse_ifcfg_file_eth_list
    global ifcfg_eth_prefix
    find_eth = {}
    for i in unuse_ifcfg_file_eth_list:
        device = ''
        master = ''
        onboot = ''
        slave = ''
        startmode = ''
        try:
            conf = configobj.ConfigObj(ifcfg_eth_prefix+i)
        except:
            ##sign = 1
            print "[INFO]ifcfg_has_dumplicated_key_or_format_err",ifcfg_eth_prefix+i
            try:
                os.system("sudo cat %s |tr -d ' ' |sort|uniq >/tmp/ifcfg-%s" %(ifcfg_eth_prefix+i,i))
                conf = configobj.ConfigObj("/tmp/ifcfg-"+i)
                print "[INFO]ifcfg_has_dumplicated_key",ifcfg_eth_prefix+i
                reduce_ifcfg_file_sameline("ifcfg-"+i)
            except:
                continue
        try:
            device = conf["DEVICE"]
            master = conf["MASTER"]
            onboot = conf["ONBOOT"]
            slave = conf["SLAVE"]
        except:
            try:
                device = conf["DEVICE"]
                master = conf["MASTER"]
                startmode = conf["STARTMODE"]
                slave = conf["SLAVE"]
            except:
                continue
        if device == eth and master == bond and slave == "yes":
            if onboot == "yes" or startmode == 'auto':
                find_eth[i]='on'
            else:
                find_eth[i]='off'
        
    return find_eth
    
def check_eth_file(eth,bond):
    global ifcfg_eth_prefix
    global sign 
    global unuse_ifcfg_file_eth_list
    device = ''
    master = ''
    onboot = ''
    slave = ''
    startmode = ''
    content = commands.getoutput("sudo cat {}".format(ifcfg_eth_prefix+eth))
    try:
        conf = configobj.ConfigObj(ifcfg_eth_prefix+eth)
    except:
        ##sign = 1
        print "[INFO]ifcfg_has_dumplicated_key_or_format_err",ifcfg_eth_prefix+eth
        try:
            os.system("sudo cat %s |tr -d ' ' |sort|uniq >/tmp/ifcfg-%s" %(ifcfg_eth_prefix+eth,eth))
            conf = configobj.ConfigObj("/tmp/ifcfg-"+eth)
            print "[INFO]ifcfg_has_dumplicated_key",ifcfg_eth_prefix+eth
            reduce_ifcfg_file_sameline("ifcfg-"+eth)
        except:
            find_eth = from_unuse_ifcfg_file_find_eth_config(eth,bond)
            if len(find_eth.keys()) == 0:
                sign = 1
                print "[ERROR]file_open_err",ifcfg_eth_prefix+eth,"\n",content
            elif len(find_eth.keys()) == 1:
                unuse_ifcfg_file_eth_list.append(eth)
                ##sign = 1
                for i in find_eth.keys():
                    unuse_ifcfg_file_eth_list.remove(i)
                    print "[INFO]file_open_err",ifcfg_eth_prefix+eth,"but_exists_1_file_fit_config_and:",ifcfg_eth_prefix+i,"and_onboot_or_startmode:",find_eth[i]
                    if find_eth[i] == "off":
                        #sign = 1
                        cmd1 = "sudo sed -i \"/ONBOOT/d\" {}".format(ifcfg_eth_prefix+i)
                        cmd2 = "sudo sed -i \"/STARTMODE/d\" {}".format(ifcfg_eth_prefix+i)
                        cmd3 = "sudo sed -i \"/DEVICE/aONBOOT=yes\" {}".format(ifcfg_eth_prefix+i)
                        print "[FIX]now_fix:",cmd1,cmd2,cmd3
                        ret1,output1 = commands.getstatusoutput(cmd1)
                        ret2,output2 = commands.getstatusoutput(cmd2)
                        ret3,output3 = commands.getstatusoutput(cmd3)
                        if ret1 != 0 or ret2 != 0 or ret3 != 0:
                            sign = 1 
            else:
                sign = 1
                unuse_ifcfg_file_eth_list.append(eth)
                for i in find_eth.keys():
                    unuse_ifcfg_file_eth_list.remove(i)
                    print "[ERROR]file_open_err",ifcfg_eth_prefix+eth,"but_exists_"+str(len(find_eth.keys()))+"_file_fit_config_and_one:",ifcfg_eth_prefix+i,"and_onboot_or_startmode:",find_eth[i]
            return
    try:
        device = conf["DEVICE"]
        master = conf["MASTER"]
        onboot = conf["ONBOOT"]
        slave = conf["SLAVE"]
    except:
        try:
            device = conf["DEVICE"]
            master = conf["MASTER"]
            startmode = conf["STARTMODE"]
            slave = conf["SLAVE"]
        except:
            find_eth = from_unuse_ifcfg_file_find_eth_config(eth,bond)
            if len(find_eth.keys()) == 0:
                sign = 1
                print "[ERROR]file_read_err",ifcfg_eth_prefix+eth,"\n",content
            elif len(find_eth.keys()) == 1:
                ##sign = 1
                unuse_ifcfg_file_eth_list.append(eth)
                for i in find_eth.keys():
                    unuse_ifcfg_file_eth_list.remove(i)
                    print "[INFO]file_read_err",ifcfg_eth_prefix+eth,"but_exists_1_file_fit_config_and:",ifcfg_eth_prefix+i,"and_onboot_or_startmode:",find_eth[i]
                    if find_eth[i] == "off":
                        #sign = 1
                        cmd1 = "sudo sed -i \"/ONBOOT/d\" {}".format(ifcfg_eth_prefix+i)
                        cmd2 = "sudo sed -i \"/STARTMODE/d\" {}".format(ifcfg_eth_prefix+i)
                        cmd3 = "sudo sed -i \"/DEVICE/aONBOOT=yes\" {}".format(ifcfg_eth_prefix+i)
                        print "[FIX]now_fix:",cmd1,cmd2,cmd3
                        ret1,output1 = commands.getstatusoutput(cmd1)
                        ret2,output2 = commands.getstatusoutput(cmd2)
                        ret3,output3 = commands.getstatusoutput(cmd3)
                        if ret1 != 0 or ret2 != 0 or ret3 != 0:
                            sign = 1
            else:
                sign = 1
                unuse_ifcfg_file_eth_list.append(eth)
                for i in find_eth.keys():
                    unuse_ifcfg_file_eth_list.remove(i)
                    print "[ERROR]file_read_err",ifcfg_eth_prefix+eth,"but_exists_"+str(len(find_eth.keys()))+"_file_fit_config_and_one:",ifcfg_eth_prefix+i,"and_onboot_or_startmode:",find_eth[i]
            return
    if device == eth and master == bond and slave == "yes" :
        if onboot == "yes" or startmode == 'auto':
            return
        else:
            print "[INFO]file_onboot_or_startmode_is_off",ifcfg_eth_prefix+eth,"\n",content
            #sign = 1
            cmd1 = "sudo sed -i \"/ONBOOT/d\" {}".format(ifcfg_eth_prefix+eth)
            cmd2 = "sudo sed -i \"/STARTMODE/d\" {}".format(ifcfg_eth_prefix+eth)
            cmd3 = "sudo sed -i \"/DEVICE/aONBOOT=yes\" {}".format(ifcfg_eth_prefix+eth)
            print "[FIX]now_fix:",cmd1,cmd2,cmd3
            ret1,output1 = commands.getstatusoutput(cmd1)
            ret2,output2 = commands.getstatusoutput(cmd2)
            ret3,output3 = commands.getstatusoutput(cmd3)
            if ret1 != 0 or ret2 != 0 or ret3 != 0:
                sign = 1
            return
    else:
        print "[ERROR]file_content_err",ifcfg_eth_prefix+eth,"\n",content
        sign = 1
        return

if __name__ == '__main__':
    #check bond ifcfg file is exsists
    bond_eth_list = get_bond_eth_list()
    check_ifexists_ifcfg_bond(bond_eth_list)
    #fix ifcfg file same line and special word
    for ifcfg_file in get_and_bak_ifcfg_file():
        reduce_ifcfg_file_sameline(ifcfg_file)
        fix_ifcfg_special_word(ifcfg_file)
    #fix ipv6 addr
    for ip,eth in get_ip().items():
        ipv6 = get_ipv6(ip)
        if ipv6 != "":
            fix_ifcfg_file_ipv6(eth,ip,ipv6)
    #print bond_eth_list
    check_unuse_ifcfg_bond(bond_eth_list)
    bond_id_list = get_bond_list()
    slave_for_bond = {}
    for bond_id in bond_id_list:
        slave_for_bond.update(get_bond_slave(bond_id))
    need_ifcfg_eth = slave_for_bond.keys()+bond_eth_list
    print "[INFO]实际使用需要ifcfg文件的网卡:",json.dumps(need_ifcfg_eth).replace(', ','+')
    #check if slave eth for bond lost ifcfg_file
    unuse_ifcfg_file_eth_list = check_unuse_ifcfg_interface(slave_for_bond)
    lose_ifcfg_eth_list = check_lost_ifcfg_interface(slave_for_bond)
    for i in slave_for_bond.keys():
        if i in lose_ifcfg_eth_list:
            find_eth = from_unuse_ifcfg_file_find_eth_config(i,slave_for_bond[i])
            if len(find_eth.keys()) == 0:
                sign = 1
                print "[ERROR]not_exists:",ifcfg_eth_prefix+i,"belongs_to:",slave_for_bond[i]
            elif len(find_eth.keys()) == 1:
                ##sign = 1
                for eth in find_eth.keys():
                    unuse_ifcfg_file_eth_list.remove(eth)
                    print "[INFO]not_exists:",ifcfg_eth_prefix+i,"but_exists_1_file_fit_config_and:",ifcfg_eth_prefix+eth,"and_onboot_or_startmode:",find_eth[eth]
                    if find_eth[eth] == "off":
                        cmd1 = "sudo sed -i \"/ONBOOT/d\" {}".format(ifcfg_eth_prefix+eth)
                        cmd2 = "sudo sed -i \"/STARTMODE/d\" {}".format(ifcfg_eth_prefix+eth)
                        cmd3 = "sudo sed -i \"/DEVICE/aONBOOT=yes\" {}".format(ifcfg_eth_prefix+eth)
                        print "[FIX]now_fix:",cmd1,cmd2,cmd3
                        ret1,output1 = commands.getstatusoutput(cmd1)
                        ret2,output2 = commands.getstatusoutput(cmd2)
                        ret3,output3 = commands.getstatusoutput(cmd3)
                        if ret1 != 0 or ret2 != 0 or ret3 != 0:
                            sign = 1
                        #sign = 1
            else:
                sign = 1
                for eth in find_eth.keys():
                    unuse_ifcfg_file_eth_list.remove(eth)
                    print "[ERROR]not_exists:",ifcfg_eth_prefix+i,"but_exists_"+str(len(find_eth.keys()))+"_file_fit_config_and_one:",ifcfg_eth_prefix+eth,"and_onboot_or_startmode:",find_eth[eth]
        else:
            check_eth_file(i,slave_for_bond[i])

    #check if exists no-use ifcfg_file
    if len(unuse_ifcfg_file_eth_list) != 0:
        os.system("sudo mkdir -p {}".format(ifcfg_path+"ifcfgfixbak/"))
        for i in unuse_ifcfg_file_eth_list:
            #ret,output = commands.getstatusoutput("sudo cat %s |tr -d ' '|grep -v ^#|grep -i onboot|grep -i yes" %(ifcfg_eth_prefix+i))
            #ret1,output1 = commands.getstatusoutput("sudo cat %s |tr -d ' '|grep -v ^#|grep -i startmode|grep -i auto" %(ifcfg_eth_prefix+i))
            #if ret == 0 or ret1 == 0:
            #    #sign = 1
            #    print "unuse_ifcfg_eth_file:",ifcfg_eth_prefix+i,"is onboot=yes and startmode=auto!"
            #else:
            #    print "unuse_ifcfg_eth_file:",ifcfg_eth_prefix+i
            #sign = 1
            print "[INFO]unuse_ifcfg_eth_file:",ifcfg_eth_prefix+i
            cmd = "sudo mv {} {}".format(ifcfg_eth_prefix+i,ifcfg_path+"ifcfgfixbak/")
            print "[FIX]now_fix:",cmd
            ret,output = commands.getstatusoutput(cmd)
            if ret != 0:
                sign = 1
    if len(unuse_ifcfg_file_bond_list) != 0:
        os.system("sudo mkdir -p {}".format(ifcfg_path+"ifcfgfixbak/"))
        for i in unuse_ifcfg_file_bond_list:
            #sign = 1
            print "[INFO]unuse_ifcfg_bond_file:",ifcfg_eth_prefix+i
            cmd = "sudo mv {} {}".format(ifcfg_eth_prefix+i,ifcfg_path+"ifcfgfixbak/")
            print "[FIX]now_fix:",cmd
            ret,output = commands.getstatusoutput(cmd)
            if ret != 0:
                sign = 1
    if sign == 1:
        exit(1)