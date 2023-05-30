#!/usr/bin/python
import commands
import json
import os
from xml.dom.minidom import parseString

DBADDRESS = "/home/kylin-ksvd/.ksvd-local/dbaddress"
SERVERSETTINGS = "/var/lib/ksvd/settings.node"


def getServerIp():
    with open(SERVERSETTINGS) as fp:
        lines = fp.readlines()
    for i in range(0, len(lines)):
        if "UNIQB_KSVD_PUBADDR" in lines[i]:
            pubAddr=lines[i].strip('\n').split('=')[1].strip('"')
            break
    return pubAddr


# get master mc address
def getMasterMcIp():
    # get master mc ip
    masterMc=''
    if not os.path.exists(DBADDRESS):
        return masterMc

    with open(DBADDRESS) as fp:
        lines = fp.readlines()
        for i in range(0, len(lines)):
            if "UNIQB_KSVD_CMADDR" in lines[i]:
                masterMc = lines[i].strip('\n').split('=')[1].strip('"')
    return masterMc


# get storage info
def getStorageInfo():
    cmd = "/usr/lib/ksvd/bin/ksvd-local-cli --action='list-data-storage'|jq '.stdout'|tr -d \" \n\""
    rc, dataInfo = commands.getstatusoutput(cmd)
    storageInfoDict = {}
    if 0 == rc:
        dataDict = json.loads(dataInfo)
        if len(dataDict) != 0:
            for key in dataDict:
                storageInfoDict[key] = dataDict[key]['type']
    cmd = "/usr/lib/ksvd/bin/ksvd-local-cli --action='list-manage-storage'|jq '.stdout'"
    rc, manageInfo = commands.getstatusoutput(cmd)
    if 0 == rc:
        manageDict = json.loads(manageInfo)
        if len(manageDict) != 0:
            storageInfoDict['MStorage'] = manageDict['type']
    return storageInfoDict


# create information from xml string
def parseXmlData(str, bricks):
    domTree = parseString(str)
    rootNode = domTree.documentElement
    nodes = rootNode.getElementsByTagName('node')
    for node in nodes:
        statusNode = node.getElementsByTagName("status")[0]
        brickPathNode = node.getElementsByTagName("path")[0]
        brickNameNode = node.getElementsByTagName("hostname")[0]
        brickName = brickNameNode.childNodes[0].nodeValue + ":" + brickPathNode.childNodes[0].nodeValue
        if brickName in bricks:
            bricks.remove(brickName)
            if statusNode.childNodes[0].nodeValue != "1":
                return 1
    return 0


# collect gluster storage brick information record when brick is offline
def collectBrickOffline():
    glusterWarnList = []

    localIp = getServerIp()
    # only master mc node report storage warning
    masterMc = getMasterMcIp()
    if localIp != masterMc:
        return glusterWarnList
    storageInfos = getStorageInfo()
    for key in storageInfos:
        if storageInfos[key] == "GlusterFS" and key != "MStorage":
            # data gluster
            cmd = "gluster volume info " + key + "|egrep 'Brick.*node'|awk -F ' ' '{print $2}'"
            rc, result = commands.getstatusoutput(cmd)
            if 0 == rc:
                bricks = result.strip('\n').split()
            if len(bricks) == 0:
                return glusterWarnList
            rc, dataGlusterStatusXml = commands.getstatusoutput("gluster volume status " + key + " --xml")
            if 0 == rc:
                dataGlusterStatus = parseXmlData(dataGlusterStatusXml, bricks)
                if dataGlusterStatus == 1:
                    glusterWarnList.append(key)
        if storageInfos[key] == "GlusterFS" and key == "MStorage":
            # manage gluster
            cmd = "/usr/share/cockpit/virtualization/bin/cockpit_cli.sh --action get_node_info --field ip --is_docker_node 1"
            rc, hostIps = commands.getstatusoutput(cmd)
            if 0 == rc:
                hostIpList = hostIps.strip('\n').split()
                for host_ip in hostIpList:
                    cmd = "ssh root@" + host_ip + " exec 'docker exec ksvd-MStorage gluster volume info MStorage'|egrep 'Brick.*node'|awk -F ' ' '{print $2}'"
                    rc, result = commands.getstatusoutput(cmd)
                    if 0 == rc:
                        bricks = result.strip('\n').split()
                    if len(bricks) == 0:
                        continue
                    cmd = "ssh root@" + host_ip + " exec 'docker exec ksvd-MStorage gluster volume status MStorage --xml'"
                    rc, manageGlusterStatusXml = commands.getstatusoutput(cmd)
                    if 0 == rc:
                        manageGlusterStatus = parseXmlData(manageGlusterStatusXml, bricks)
                        if manageGlusterStatus == 1:
                            glusterWarnList.append(key)
                        break
    return glusterWarnList


def main():
    offlineWarn = collectBrickOffline()
    print(json.dumps(offlineWarn))


if __name__ == "__main__":
    main()

