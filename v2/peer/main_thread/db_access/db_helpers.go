package db_helpers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"

	ftp_context "github.com/it-shiloheye/ftp_system/v2/lib/context"
	db "github.com/it-shiloheye/ftp_system/v2/lib/db_access"
	db_access "github.com/it-shiloheye/ftp_system/v2/lib/db_access/generated"
	"github.com/it-shiloheye/ftp_system/v2/lib/logging"
	"github.com/it-shiloheye/ftp_system/v2/lib/logging/log_item"
	ftp_tlshandler "github.com/it-shiloheye/ftp_system/v2/lib/tls_handler/v2"
	server_config "github.com/it-shiloheye/ftp_system/v2/peer/config"
)

var DB = db.DB
var ServerConfig = server_config.ServerConfig

var Logger = logging.Logger

func ticker(loc log_item.Loc, i int, v ...any) {
	log.Printf("%s\n%03d: %s\n", string(loc), i, fmt.Sprint(v...))
}

func ConnectClient(ctx ftp_context.Context, StorageStruct *server_config.StorageStruct) error {
	loc := log_item.Loc(`ConnectClient(ctx ftp_context.Context) error`)
	var err1, err3, err4, err5, err6, err7 error
	var db_peers []*db_access.PeersTable
	if db.DBPool.Len() < 10 {
		// log.Println("connect client checking length")
		db.DBPool.PopulateConns(ctx.Add(), 10)
		ctx.Finished()
	}

	db_conn := db.DBPool.GetConn()
	defer db.DBPool.Return(db_conn)

	ip_addr := ServerConfig.LocalIp()

	// ticker(loc, 1, "before query")
	db_peers, err1 = DB.ConnectClient(context.TODO(), db_conn, &db_access.ConnectClientParams{
		IpAddress: ip_addr.String(),
		PeerID:    StorageStruct.PeerId,
	})
	if err1 != nil {
		return Logger.LogErr(loc, err1)
	}

	// ticker(loc, 2, "before template_cd")
	template_cd := ftp_tlshandler.CertData{
		Organization:  "Shiloh Eye, Ltd",
		Country:       "KE",
		Province:      "Coast",
		Locality:      "Mombasa",
		StreetAddress: "2nd Floor, SBM Bank Centre, Nyerere Avenue, Mombasa",
		PostalCode:    "80100",
		NotAfter: ftp_tlshandler.NotAfterStruct{
			Days: 31,
		},
		IPAddrresses: []net.IP{
			net.IPv4(127, 0, 0, 1),
			net.IPv6loopback,
			ip_addr.To4(),
		},
	}

	if len(db_peers) > 0 {
		// ticker(loc, 3, "need to decode PEM")
		db_peer := db_peers[0]
		StorageStruct.PeerId = db_peer.PeerID
		tmp_ca_pem := ftp_tlshandler.CAPem{}

		err2 := json.Unmarshal(db_peer.Pem, &tmp_ca_pem)
		if err2 != nil {
			log.Println(err2)
			return Logger.LogErr(loc, err2)
		}

		// generate new tls each time
		x509_tls_cert := ftp_tlshandler.ExampleTLSCert(template_cd)

		ServerConfig.TLS_Cert, err3 = ftp_tlshandler.GenerateTLSCert(&tmp_ca_pem, x509_tls_cert)
		if err3 != nil {
			return Logger.LogErr(loc, err3)
		}

		log.Println("succeffully loaded pem from database")
		return nil
	}
	log.Println("registering new client to db")

	// simple time guard, update cert every 7 days, server restarts every day at least once
	ServerConfig.TLS_Cert_Creation = time.Now()

	tmp_x509 := ftp_tlshandler.ExampleCACert(template_cd)

	tmp_ca_pem, err4 := ftp_tlshandler.GenerateCAPem(tmp_x509)
	if err4 != nil {
		return Logger.LogErr(loc, err4)
	}

	// generate new tls each time
	x509_tls_cert := ftp_tlshandler.ExampleTLSCert(template_cd)

	// ticker(loc, 4, "before tls")
	ServerConfig.TLS_Cert, err5 = ftp_tlshandler.GenerateTLSCert(&tmp_ca_pem, x509_tls_cert)
	if err5 != nil {
		return Logger.LogErr(loc, err5)
	}
	// ticker(loc, 5, "before encoding to ca_pem")

	var d []byte
	d, err6 = json.Marshal(&tmp_ca_pem)
	if err6 != nil {
		return Logger.LogErr(loc, err6)
	}

	tmp_peerids, err7 := DB.CreateClient(ctx, db_conn, &db_access.CreateClientParams{
		IpAddress: ip_addr.String(),
		Pem:       d,
		PeerRole:  db_access.NullPeerRoleType{PeerRoleType: StorageStruct.PeerRole, Valid: true},
	})
	if err7 != nil {
		return Logger.LogErr(loc, err7)
	}

	StorageStruct.PeerId = tmp_peerids[0].PeerID
	StorageStruct.PeerRole = tmp_peerids[0].PeerRole.PeerRoleType

	return nil
}

func GetFiles(ctx ftp_context.Context, StorageStruct *server_config.StorageStruct) ([]*db_access.GetFilesRow, error) {
	loc := log_item.Locf(`func GetFiles(ctx ftp_context.Context, StorageStruct: %v) error `, StorageStruct.PeerId.Bytes)
	conn := db.DBPool.GetConn()
	defer db.DBPool.Return(conn)

	log.Println("before DB.GetFiles")
	defer log.Println("after DB.GetFiles")
	f, err1 := DB.GetFiles(ctx, conn, StorageStruct.PeerId.Bytes)
	if err1 != nil {

		return nil, Logger.LogErr(loc, err1)
	}

	return f, nil
}

func UpdatePeerRole(ctx ftp_context.Context, StorageStruct *server_config.StorageStruct, prev_role *db_access.PeerRoleType) error {
	loc := log_item.Locf(`func UpdatePeerRole(ctx ftp_context.Context, StorageStruct: %v, prev_role: %v) error`, StorageStruct.PeerId.Bytes, prev_role)
	conn := db.DBPool.GetConn()
	defer db.DBPool.Return(conn)

	if *prev_role == StorageStruct.PeerRole {
		return nil
	}

	err1 := DB.UpdatePeerRole(ctx, conn, &db_access.UpdatePeerRoleParams{
		PeerID:   StorageStruct.PeerId,
		PeerRole: db_access.NullPeerRoleType{PeerRoleType: StorageStruct.PeerRole, Valid: true},
	})

	if err1 != nil {
		Logger.LogErr(loc, err1)
	}

	return nil
}

func MarkDeleted(ctx ftp_context.Context, StorageStruct *server_config.StorageStruct, file_hash string) (bool, error) {
	loc := log_item.Locf(`MarkDeleted(ctx ftp_context.Context, StorageStruct: %v, file_hash: %v) (bool, error)`, StorageStruct.PeerId.Bytes, file_hash)
	conn := db.DBPool.GetConn()
	defer db.DBPool.Return(conn)

	stored_count, err1 := DB.CountIfStored(ctx, conn, &file_hash)
	if err1 != nil {
		return false, Logger.LogErr(loc, err1)
	}

	to_delete := stored_count > 0
	if !to_delete {
		return false, nil
	}

	err2 := DB.MarkFileDeleted(ctx, conn, &db_access.MarkFileDeletedParams{
		FileHash: &file_hash,
		PeerID:   StorageStruct.PeerId.Bytes,
	})

	if err2 != nil {
		return false, Logger.LogErr(loc, err2)
	}

	return to_delete, nil
}
