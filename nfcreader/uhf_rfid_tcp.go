package uhfrftcp

import (
	"flag"
	"fmt"
	"io"
	"log"
	apihttp "mantap2/api/http"
	"mantap2/config"
	"net"
	"os"
	"strings"
	"time"
)

var last_sended_sn_rfid = ""

func wrLog(isdevonly bool, msg string) {
	if config.Log_show {
		if isdevonly {
			if config.Log_dev {
				log.Print(msg)
			}
		} else {
			log.Print(msg)
		}
	}
}

// Handles TC connection and perform synchorinization:
// TCP -> Stdout and Stdin -> TCP
func tcp_con_handle(con net.Conn) {
	wrLog(true, "UHF Reader tcp_con_handle")
	chan_to_stdout := stream_copy(con, os.Stdout)

	chan_to_remote := stream_copy(os.Stdin, con)
	select {
	case <-chan_to_stdout:
		wrLog(false, "UHF Reader Remote connection is closed")
	case <-chan_to_remote:
		wrLog(false, "UHF Reader Local program is terminated")
	}
}

// Performs copy operation between streams: os and tcp streams
func stream_copy(src io.Reader, dst io.Writer) <-chan int {
	wrLog(true, "UHF Reader stream_copy")
	buf := make([]byte, 1024)
	sync_channel := make(chan int)
	go func() {
		defer func() {
			if con, ok := dst.(net.Conn); ok {
				con.Close()
				log.Printf("UHF Reader Connection from %v is closed\n", con.RemoteAddr())
			}
			sync_channel <- 0 // Notify that processing is finished
		}()
		for {

			var nBytes int
			var err error
			nBytes, err = src.Read(buf)
			if err != nil {
				if err != io.EOF {
					log.Printf("UHF Reader Read error: %s\n", err)
				}
				break
			}

			// _, err = dst.Write(buf[0:nBytes])
			// if err != nil {
			// 	log.Fatalf("Write error: %s\n", err)
			// }

			//log.Println(buf[0:nBytes])
			sHexnya := ""
			for j := 0; j < nBytes; j++ {
				h := ""
				if buf[j] == 0 {
					h = "00"
				} else {
					h = fmt.Sprintf("%x", buf[j])
				}
				sHexnya = strings.TrimSpace(sHexnya + " " + h)

			}
			if last_sended_sn_rfid != sHexnya {
				log.Println(sHexnya)
				//sendToHttpAPI("JTI-Gate01", sHexnya)
				apihttp.SendToHttpAPI(config.Vr_gen_nama, sHexnya)
			}

			last_sended_sn_rfid = sHexnya

		}
	}()
	return sync_channel
}

// Accept data from UPD connection and copy it to the stream
func accept_from_udp_to_stream(src net.Conn, dst io.Writer) <-chan net.Addr {
	buf := make([]byte, 1024)
	sync_channel := make(chan net.Addr)
	con, ok := src.(*net.UDPConn)
	if !ok {
		log.Printf("UHF Reader Input must be UDP connection")
		return sync_channel
	}
	go func() {
		var remoteAddr net.Addr
		for {
			var nBytes int
			var err error
			var addr net.Addr
			nBytes, addr, err = con.ReadFromUDP(buf)
			if err != nil {
				if err != io.EOF {
					log.Printf("UHF Reader Read error: %s\n", err)
				}
				break
			}
			if remoteAddr == nil && remoteAddr != addr {
				remoteAddr = addr
				sync_channel <- remoteAddr
			}
			_, err = dst.Write(buf[0:nBytes])
			if err != nil {
				log.Fatalf("UHF Reader Write error: %s\n", err)
			}
		}
	}()
	log.Println("UHF Reader Exit write_from_udp_to_stream")
	return sync_channel
}

// Put input date from the stream to UDP connection
func put_from_stream_to_udp(src io.Reader, dst net.Conn, remoteAddr net.Addr) <-chan net.Addr {
	buf := make([]byte, 1024)
	sync_channel := make(chan net.Addr)
	go func() {
		for {
			var nBytes int
			var err error
			nBytes, err = src.Read(buf)
			if err != nil {
				if err != io.EOF {
					log.Printf("UHF Reader Read error: %s\n", err)
				}
				break
			}
			log.Println("UHF Reader Write to the remote address:", remoteAddr)
			if con, ok := dst.(*net.UDPConn); ok && remoteAddr != nil {
				_, err = con.WriteTo(buf[0:nBytes], remoteAddr)
			}
			if err != nil {
				log.Fatalf("UHF Reader Write error: %s\n", err)
			}
		}
	}()
	return sync_channel
}

// Handle UDP connection
func udp_con_handle(con net.Conn) {
	in_channel := accept_from_udp_to_stream(con, os.Stdout)
	log.Println("UHF Reader Waiting for remote connection")
	remoteAddr := <-in_channel
	log.Println("UHF Reader Connected from", remoteAddr)
	out_channel := put_from_stream_to_udp(os.Stdin, con, remoteAddr)
	select {
	case <-in_channel:
		log.Println("UHF Reader Remote connection is closed")
	case <-out_channel:
		log.Println("UHF Reader Local program is terminated")
	}
}

func DoJobs() {

	wrLog(false, "UHF Reader Trying Connect to "+config.Rdr_ip+" PORT "+config.Rdr_port)

	if !config.IsUdp {
		wrLog(false, "UHF Reader Work with TCP protocol")
		if config.IsListen {
			listener, err := net.Listen("tcp", config.Rdr_port)
			if err != nil {
				log.Fatalln(err)
			}
			log.Println("UHF Reader Listening on", config.Rdr_port)
			con, err := listener.Accept()
			if err != nil {
				log.Fatalln(err)
			}
			log.Println("UHF Reader Connect from", con.RemoteAddr())
			tcp_con_handle(con)

		} else if config.Rdr_ip != "" {
			for {
				//con, err := net.Dial("tcp", rdr_ip+":"+rdr_port)
				con, err := net.DialTimeout("tcp", config.Rdr_ip+":"+config.Rdr_port, config.Rdr_timeout)
				//net.TCPConn.SetKeepAlive/net.TCPConn.SetKeepAlivePeriod/net.TCPConn.SetNoDelay
				//conn.SetKeepAlive(true)
				//conn.SetKeepAlivePeriod(time.Second * 60)

				if err != nil {
					log.Println(err)
					time.Sleep(3 * time.Second)
				} else {
					log.Println("UHF Reader Connected to", config.Rdr_ip+":"+config.Rdr_port)
					tcp_con_handle(con)
				}
				wrLog(true, "-> UHF Reader Loooping")
			}
		} else {
			flag.Usage()
		}
	} else {
		wrLog(false, "UHF Reader Work with UDP protocol")
		if config.IsListen {
			addr, err := net.ResolveUDPAddr("udp", config.Rdr_port)
			if err != nil {
				log.Fatalln(err)
			}
			con, err := net.ListenUDP("udp", addr)
			if err != nil {
				log.Fatalln(err)
			}
			log.Println("UHF Reader Has been resolved UDP address:", addr)
			log.Println("UHF Reader Listening on", config.Rdr_port)
			udp_con_handle(con)
		} else if config.Rdr_ip != "" {
			addr, err := net.ResolveUDPAddr("udp", config.Rdr_ip+":"+config.Rdr_port)
			if err != nil {
				log.Fatalln(err)
			}
			log.Println("UHF Reader Has been resolved UDP address:", addr)
			con, err := net.DialUDP("udp", nil, addr)
			if err != nil {
				log.Fatalln(err)
			}
			udp_con_handle(con)
		}
	}
}
