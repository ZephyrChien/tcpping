package main
import(
	"os"
	"fmt"
	"net"
	"time"
	"flag"
)
var(
	host=flag.String("h","github.com","ip or domain")
	port=flag.Int("p",443,"port")
	count=flag.Int("n",4,"count")
	timeout=flag.Int("t",1,"timeout of each request")
	quiet=flag.Bool("q",false,"work in quiet mode")
)
func main(){
	flag.Parse()
	ip:=resolv(*host)
	target:=fmt.Sprintf("%s:%d",ip,*port)
	results:=make([]int,*count)
	fmt.Printf("TCPPING %s (%s):\n",*host,ip)
	for i:=0;i<*count;i++{
		results[i]=tcping(target,i+1)
	}
	analysis(results[:])
}
func resolv(address string)net.IP{
	ips,err:=net.LookupIP(address)
	if err!=nil{
		fmt.Fprintln(os.Stderr,err)
	}
	return ips[0]
}
func tcping(target string,seq int)int{
	start:=time.Now()
	conn,err:=net.DialTimeout("tcp",target,time.Duration(*timeout)*time.Second)
	end:=time.Now()
	latency:=int(end.Sub(start).Milliseconds())
	var status string
	if err!=nil{
		status="close"
	}else{
		conn.Close()
		status="open"
	}
	if !*quiet{
		fmt.Printf("seq%4d: %s[%s] %dms\n",seq,target,status,latency)
	}
	return latency
}
func analysis(results []int){
	min,max,sum:=results[0],results[0],0
	for _,val:=range results{
		sum += val
		if val < min{
			min=val
		}
		if val > max{
			max=val
		}
	}
	fmt.Printf("min/avg/max = %d/%d/%dms\n",min,sum/(*count),max)
}
