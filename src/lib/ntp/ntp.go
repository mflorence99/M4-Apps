package ntp

import (
	_ "embed"
	"errors"
	"io"
	"math/rand"
	"net"
	"runtime"
	"time"
)

// 👁️ https://github.com/tinygo-org/drivers/blob/release/examples/net/ntpclient/main.go

const NTP_PACKET_SIZE = 48

// 👇 70 years is diff between 1900 (NTP time) and 1970 (Unix time)
const NTP2UNIXTIME = 2208988800

var servers = []string{
	"time.cloudflare.com",
	"time1.google.com",
	"time2.google.com",
	"time3.google.com",
	"time4.google.com",
	"0.pool.ntp.org",
	"1.pool.ntp.org",
	"2.pool.ntp.org",
	"3.pool.ntp.org",
}

// 🟧 Synchronize system time with NTP

func SyncSystemTime() error {

	// 👇 pick a random server
	ix := rand.Intn(len(servers))
	server := servers[ix] + ":123"
	println("🐞 requesting NTP time from", server)

	// 👇 make a UDP connection
	conn, err := net.Dial("udp", server)
	if err != nil {
		return (err)
	}

	// 👇 ... and close it when we're done
	defer conn.Close()

	// 👇 carefully parse its reply as the current time
	t, err := getCurrentTime(conn)
	if err != nil {
		return err
	}

	// 👇 finally!
	runtime.AdjustTimeOffset(-1 * int64(time.Since(t)))
	return nil
}

// 🟦 Helpers

func getCurrentTime(conn net.Conn) (time.Time, error) {

	// 👇 send request
	var request = [NTP_PACKET_SIZE]byte{0xe3}
	if _, err := conn.Write(request[:]); err != nil {
		return time.Time{}, err
	}

	// 👇 extract response
	var response = make([]byte, NTP_PACKET_SIZE)
	n, err := conn.Read(response)
	if err != nil && err != io.EOF {
		return time.Time{}, err
	}
	if n != NTP_PACKET_SIZE {
		return time.Time{}, errors.New("unexpected NTP packet size")
	}

	// 👇 the timestamp starts at byte 40 of the received packet
	//    and is four bytes, this is NTP time (seconds since Jan 1 1900)
	t := uint32(response[40])<<24 | uint32(response[41])<<16 | uint32(response[42])<<8 | uint32(response[43])
	return time.Unix(int64(t-NTP2UNIXTIME), 0), nil
}
