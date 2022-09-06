for i in `cat /proc/partitions | awk {'print $4'} | grep sd`
do
echo "Device: $i WWID: `scsi_id -g -u -s /block/$i`"
done | sort -k4

