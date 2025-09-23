host=$(ip route show | grep -i default | awk '{ print $3}')
echo $host
