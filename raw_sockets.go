package main

import (
	"fmt"
	"os"
	"net"
	//"log"
	//"syscall"
	//"bufio"
	"bytes"
	"encoding/binary"
	"flag"
	"runtime"
	"golang.org/x/sys/unix"
)

const maxPacketSize = 2 << 16

type iphdr struct {
	vhl   uint8
	tos   uint8
	iplen uint16
	id    uint16
	off   uint16
	ttl   uint8
	proto uint8
	csum  uint16
	src   [4]byte
	dst   [4]byte
}

type udphdr struct {
	src  uint16
	dst  uint16
	ulen uint16
	csum uint16
}

// pseudo header used for checksum calculation
type pseudohdr struct {
	ipsrc   [4]byte
	ipdst   [4]byte
	zero    uint8
	ipproto uint8
	plen    uint16
}

func checksum(buf []byte) uint16 {

	sum := uint32(0)

	for ;len(buf) >=2;buf=buf[2:]{
		sum = sum + uint32(buf[0] << 8) | uint32((buf[1]))
	}
	if len(buf)>0 {
		sum = sum + uint32(buf[0]) << 8
	}
	for sum > 0xffff {
		sum = (sum >> 16) + (sum & 0xffff)
	}
	csum := ^uint16(sum)


	/*
	 * From RFC 768:
	 * If the computed checksum is zero, it is transmitted as all ones (the
	 * equivalent in one's complement arithmetic). An all zero transmitted
	 * checksum value means that the transmitter generated no checksum (for
	 * debugging or for higher level protocols that don't care).
	 */
	if csum == 0 {
		csum = 0xffff
	}
	return csum
}

func (h *iphdr)checksum(){
	h.csum =0
	var b bytes.Buffer
	binary.Write(&b,binary.BigEndian,h)
	h.csum=checksum(b.Bytes())
}

func (u *udphdr)checksum(ip *iphdr, payload []byte) {
	u.csum=0
	phdr := pseudohdr{
		ipsrc:   ip.src,
		ipdst:   ip.dst,
		zero:    0,
		ipproto: ip.proto,
		plen:    u.ulen,
	}
	var b bytes.Buffer
	binary.Write(&b, binary.BigEndian, &phdr)
	binary.Write(&b, binary.BigEndian, u)
	binary.Write(&b, binary.BigEndian, &payload)
	u.csum = checksum(b.Bytes())
}

func main(){
	ipsrcaddr := "10.10.100.100"
	ipdstaddr := "10.10.10.65"

	ipsrcaddr1 := "10.10.10.65"
	ipdstaddr1 := "10.10.10.1"

	udpsrc := uint(10000)
	udpdst := uint(20000)
	flag.StringVar(&ipsrcaddr, "ipsrc", ipsrcaddr, "IPv4 source address")
	flag.StringVar(&ipdstaddr, "ipdst", ipdstaddr, "IPv4 destination address")
	flag.UintVar(&udpsrc, "udpsrc", udpsrc, "UDP source port")
	flag.UintVar(&udpdst, "udpdst", udpdst, "UDP destination port")
	flag.Parse()

	ipsrc := net.ParseIP(ipsrcaddr)
	if ipsrc == nil {
		fmt.Fprintf(os.Stderr, "invalid source IP: %v\n", ipsrc)
		os.Exit(1)
	}
	ipsrc1 := net.ParseIP(ipsrcaddr1)
	if ipsrc1 == nil {
		fmt.Fprintf(os.Stderr,"Invalid Source IP:%v\n",ipsrc1)
		os.Exit(1)
	}

	ipdst := net.ParseIP(ipdstaddr)
	if ipdst == nil {
		fmt.Fprintf(os.Stderr, "invalid destination IP: %v\n", ipdst)
		os.Exit(1)
	}
	ipdst1 := net.ParseIP(ipdstaddr1)
	if ipdst1 == nil {
		fmt.Fprintf(os.Stderr, "invalid destination IP: %v\n", ipdst)
		os.Exit(1)
	}

	fd, err := unix.Socket(unix.AF_INET, unix.SOCK_RAW, unix.IPPROTO_RAW)
	if err != nil || fd < 0 {
		fmt.Fprintf(os.Stdout, "error creating a raw socket: %v\n", err)
		os.Exit(1)
	}
	err = unix.SetsockoptInt(fd, unix.IPPROTO_IP, unix.IP_HDRINCL, 1)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error enabling IP_HDRINCL: %v\n", err)
		unix.Close(fd)
		os.Exit(1)
	}

	ip := iphdr{
		vhl: 0x45,
		tos: 0,
		// iplen set later
		id:    0,
		off:   0,
		ttl:   64,
		proto: unix.IPPROTO_UDP,
		// ipsum set later
	}

	ip1 := iphdr{
		vhl: 0x45,
		tos:0,
		id:0,
		off:0,
		ttl: 64,
		proto: unix.IPPROTO_IPIP,
	}

	copy(ip.src[:], ipsrc.To4())
	copy(ip.dst[:], ipdst.To4())

	copy(ip1.src[:],ipsrc1.To4())
	copy(ip1.dst[:],ipdst1.To4())

	udp := udphdr{
		src: uint16(udpsrc),
		dst: uint16(udpdst),
		// ulen set later
		// csum set later
	}

	addr := unix.SockaddrInet4{}
	//addr1 := unix.SockaddrInet4{}

	for {
		/*
		stdin := bufio.NewReader(os.Stdin)
		line, err := stdin.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			break
		}
		*/
		line := "This is test"
		payload := []byte(line)
		udplen := 8 + len(payload)

		totalLen := 20 + udplen
		totalLen1 := 20 + totalLen

		if totalLen > maxPacketSize {
			fmt.Fprintf(os.Stderr, "message is too large to fit into a packet: %v > %v\n", totalLen, maxPacketSize)
			continue
		}

		ip.iplen = uint16(totalLen)
		ip.checksum()

		ip1.iplen = uint16(totalLen1)
		ip1.checksum()

		udp.ulen = uint16(udplen)
		udp.checksum(&ip, payload)

		var b bytes.Buffer

		err = binary.Write(&b, binary.BigEndian, &ip1)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error encoding the IP header: %v\n", err)
			continue
		}




		err = binary.Write(&b, binary.BigEndian, &ip)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error encoding the IP header: %v\n", err)
			continue
		}
		err = binary.Write(&b, binary.BigEndian, &udp)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error encoding the UDP header: %v\n", err)
			continue
		}
		err = binary.Write(&b, binary.BigEndian, &payload)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error encoding the payload: %v\n", err)
			continue
		}
		bb := b.Bytes()
		if runtime.GOOS == "darwin" {
			bb[2], bb[3] = bb[3], bb[2]
		}
		err = unix.Sendto(fd, bb, 0, &addr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error sending the packet: %v\n", err)
			continue
		}
		fmt.Printf("%v bytes were sent\n", len(bb))

	}
	err = unix.Close(fd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error closing the socket: %v\n", err)
		os.Exit(1)
	}

}

func checkError(err error) {
    if err != nil {
        fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
    }
}
