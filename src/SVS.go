package main

import (
	"crypto/rand"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

type OMK struct {
	Total      int            `json:"total"`
	Way        int            `json:"way"`
	Pending    int            `json:"pending"`
	Accepted   int            `json:"accepted"`
	Brak       int            `json:"brak"`
	BrakItems  []ITEM         `json:"brakItems"`
	Empty      int            `json:"empty"`
	Gost       map[string]int `json:"gost"`
	EmptyItems []ITEM         `json:"emptyItems"`
}

type HISTORY struct {
	Certificate CERTIFICATE `json:"certificate"`
	Timestamp   int64       `json:"timestamp"`
	TxId        string      `json:"transactionId"`
}

type BASICSTATS struct {
	TotalLot  int `json:"totalLot"`
	TotalCert int `json:"totalCert"`
}

type SEARCH struct {
	Contracts []CONTRACT
}

type Search_Shipment struct {
	Shipments []Shipment
}

type ITEM struct {
	SUBLOT_NUMBER    string `xml:"SUBLOT_NUMBER"`
	PRODUCER_YEAR    string `xml:"PRODUCER_YEAR"`
	LOT_COAT_INT     string `xml:"LOT_COAT_INT"`
	UOM_PRESSURE     string `xml:"UOM_PRESSURE"`
	PRESSURE         string `xml:"PRESSURE"`
	LOT_WELD         string `xml:"LOT_WELD"`
	STRENGTH         string `xml:"STRENGTH"`
	LOT_COAT         string `xml:"LOT_COAT"`
	DIAM             string `xml:"DIAM"`
	BRAND            string `xml:"BRAND"`
	GOST             string `xml:"GOST"`
	GOST_COAT_EXT    string `xml:"GOST_COAT_EXT"`
	GOST_COAT_INT    string `xml:"GOST_COAT_INT"`
	SMELTING_NUMBER  string `xml:"SMELTING_NUMBER"`
	WALL             string `xml:"WALL"`
	LENGTH           string `xml:"LENGTH"`
	WEIGHT           string `xml:"WEIGHT"`
	SUBLOT_NUMBER2   string `xml:"SUBLOT_NUMBER2"`
	SMELTING_NUMBER2 string `xml:"SMELTING_NUMBER2"`
	AMOUNT_SEAM      string `xml:"AMOUNT_SEAM"`
	URL              string
	STATUS           string
	CERT_NUMBER      string
}

type ELEMENT struct {
	ELEMENT_NAME  string `xml:"ELEMENT_NAME"`
	ELEMENT_VALUE string `xml:"ELEMENT_VALUE"`
}

type SMELTING struct {
	SMELTING_NUMBER string    `xml:"SMELTING_NUMBER"`
	PLATE_PRODUCER  string    `xml:"PLATE_PRODUCER"`
	PLATE_NTD       string    `xml:"PLATE_NTD"`
	PLATE_SPEC      string    `xml:"PLATE_SPEC"`
	BRAND_STEEL     string    `xml:"BRAND_STEEL"`
	ELEMENTS        []ELEMENT `xml:"ELEMENTS>ELEMENT"`
}

type MECH_TEST struct {
	MECH_NAME  string `xml:"MECH_NAME"`
	MECH_VALUE string `xml:"MECH_VALUE"`
}

type METAL struct {
	SMELTING_NUMBER string      `xml:"SMELTING_NUMBER"`
	LOT_WELD        string      `xml:"LOT_WELD"`
	MECHANICAL      []MECH_TEST `xml:"MECHANICAL>MECH_TEST"`
}

type WELD struct {
	LOT_WELD   string      `xml:"LOT_WELD"`
	MECHANICAL []MECH_TEST `xml:"MECHANICAL>MECH_TEST"`
}

type CERTIFICATE struct {
	CERT_NUMBER      string     `xml:"CERT_NUMBER"`
	ORDER            string     `xml:"ORDER"`
	ORDER_NUMBER     string     `xml:"ORDER_NUMBER"`
	TORG_12          string     `xml:"TORG_12"`
	PACKET_NUMBER    string     `xml:"PACKET_NUMBER"`
	NAME_SHIP        string     `xml:"NAME_SHIP"`
	ADDRESS_SHIP     string     `xml:"ADDRESS_SHIP"`
	PRODUCER         string     `xml:"PRODUCER"`
	SHIP_METHOD_CODE string     `xml:"SHIP_METHOD_CODE"`
	CAR_NUMBER       string     `xml:"CAR_NUMBER"`
	WAGON_NUMBER     string     `xml:"WAGON_NUMBER"`
	DATE_SHIP        string     `xml:"DATE_SHIP"`
	URL              string     `xml:"URL"`
	ITEMS            []ITEM     `xml:"ITEMS>ITEM"`
	SMELTINGS        []SMELTING `xml:"SMELTINGS>SMELTING"`
	BASIC_METAL      []METAL    `xml:"BASIC_METAL>METAL"`
	WELDS            []WELD     `xml:"WELDS>WELD"`
	STATUS           string
}

type CONTRACT struct {
	CONTRACT_NUMBER string
	CREATION_DATE   string `json:"CREATION_DATE"`
	START_DATE      string `json:"START_DATE"`
	END_DATE        string `json:"END_DATE"`
	PURCHASER       string `json:"PURCHASER"`
	SUPPLIER        string `json:"SUPPLIER"`
	MAX_PENALTY     string `json:"MAX_PENALTY"`
	PENALTY_PER_DAY string `json:"PENALTY_PER_DAY"`
	CREATED_BY      string `json:"CREATED_BY"`
	SIGNED          string `timestamp`
	AUTHORIZED      string `timestamp`
}

type Order struct {
	AgreementID  string      `xml:"AGREEMENT_ID"`
	OrderID      string      `xml:"ORDER_ID"`
	CreatedBy    string      `xml:"CREATED_BY"`
	CreateDate   string      `xml:"CREATED_DATE"`
	CreateTime   string      `xml:"CREATED_TIME"`
	PaymentTerms string      `xml:"PAYMENT_TERMS"`
	Prepayment   string      `xml:"PREPAYMENT"`
	OrderItems   []OrderItem `xml:"ORDER_ITEMS>ORDER_ITEM"`
}

type OrderItem struct {
	OrderID    string `xml:"ORDER_ID"`
	Line       string `xml:"LINE"`
	CreatedBy  string `xml:"CREATED_BY"`
	CreateDate string `xml:"CREATED_DATE"`
	CreateTime string `xml:"CREATED_TIME"`
	ItemCode   string `xml:"ITEM_CODE"`
	ItemName   string `xml:"ITEM_NAME"`
	NetValue   string `xml:"NET_VALUE"`
	Quantity   string `xml:"QUANTITY"`
	NetPrice   string `xml:"NET_PRICE"`
	Unit       string `xml:"UNIT"`
}

type Shipment struct {
	ShipmentID     string          `xml:"SHIPMENT_ID"`
	OrderID        string          `xml:"ORDER_ID"`
	CreatedBy      string          `xml:"CREATED_BY"`
	CreateDate     string          `xml:"CREATED_DATE"`
	CreateTime     string          `xml:"CREATED_TIME"`
	ShipmentDate   string          `xml:"SHIPMENT_DATE"`
	ShipmentItems  []ShipmentItem  `xml:"SHIPMENT_ITEMS>SHIPMENT_ITEM"`
}

type ShipmentItem struct {
	ShipmentID             string `xml:"SHIPMENT_ID"`
	Line                   string `xml:"LINE"`
	CreatedBy              string `xml:"CREATED_BY"`
	CreateDate             string `xml:"CREATED_DATE"`
	CreateTime             string `xml:"CREATED_TIME"`
	ItemCode               string `xml:"ITEM_CODE"`
	ItemName               string `xml:"ITEM_NAME"`
    Lot                    string `xml:"LOT"`
    CertificationNumber    string `xml:"CERTIFICATION_NUMBER"`
    OrderId                string `xml:"ORDER_ID"`
    Position               string `xml:"POSITION"`
    CertificationDate      string `xml:"CERTIFICATION_DATE"`
    Quantity               string `xml:"QUANTITY"`    
	Unit                   string `xml:"UNIT"`
}

type Payment struct {
	PaymentID    string
	DocumentID   string
	DocumentType string
	Amount       float32
	Status       string
	Timestamp    string
}

type SVS struct {
}

var dateFormat string = "02.01.2006"
var dateDiff float64 = 5
var logger = shim.NewLogger("certificate")

func Success(rc int32, doc string, payload []byte) peer.Response {
	return peer.Response{
		Status:  rc,
		Message: doc,
		Payload: payload,
	}
}

func Error(rc int32, doc string) peer.Response {
	logger.Errorf("Error %d = %s", rc, doc)
	return peer.Response{
		Status:  rc,
		Message: doc,
	}
}

func (cc *SVS) Init(stub shim.ChaincodeStubInterface) peer.Response {
	if _, args := stub.GetFunctionAndParameters(); len(args) > 0 {
		return Error(http.StatusBadRequest, "Init: Incorrect number of arguments; no arguments were expected.")
	}
	return Success(http.StatusOK, "OK", nil)
}

func (cc *SVS) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	function, args := stub.GetFunctionAndParameters()
	switch function {
	case "createContract":
		return cc.createContract(stub, args)
	case "changeStatusContract":
		return cc.changeStatusContract(stub, args)
	case "getAllContracts":
		return cc.getAllContracts(stub, args)
	case "createOrder":
		return cc.createOrder(stub, args)
	case "getAllOrders":
		return cc.getAllOrders(stub, args)
    case "createShipment":
		return cc.createShipment(stub, args)
	case "getAllShipments":
		return cc.getAllShipments(stub, args)    

		//	case "uploadXml":
		//		return cc.uploadXml(stub, args)
		//	case "search":
		//		return cc.search(stub, args)
		//	case "searchLot":
		//		return cc.searchLot(stub, args)
		//	case "status":
		//		return cc.status(stub, args)
		//	case "basicStats":
		//		return cc.basicStats(stub, args)
		//	case "stats":
		//		return cc.stats(stub, args)
		//	case "history":
		//		return cc.history(stub, args)
	default:
		logger.Warningf("Invoke('%s') invalid!", function)
		return Error(http.StatusNotImplemented, "Invalid method!")
	}
}

// newUUID generates a random UUID according to RFC 4122
func newUUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}

func (cc *SVS) createContract(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	contractId, _ := newUUID()

	contract := &CONTRACT{CONTRACT_NUMBER: contractId, CREATION_DATE: args[0], START_DATE: args[1], END_DATE: args[2], PURCHASER: args[3], SUPPLIER: args[4], MAX_PENALTY: args[5], PENALTY_PER_DAY: args[6], CREATED_BY: args[7], SIGNED: "none", AUTHORIZED: "none"}

	if value, err := stub.GetState(contractId); err != nil || value != nil {
		return Error(401, "Already exists")
	}

	jsonDoc, _ := json.Marshal(contract)

	if err := stub.PutState(contractId, jsonDoc); err != nil {
		return Error(400, "Can't create contract")
	}
	return Success(http.StatusOK, "OK", nil)
}

func (cc *SVS) changeStatusContract(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	contract, _ := stub.GetState(args[0])
	signedType := args[1]
	signedTimestamp := args[2]
	var contractStruct CONTRACT
	json.Unmarshal(contract, &contractStruct)

	if signedType == "signed" {
		contractStruct.SIGNED = signedTimestamp
	} else if signedType == "authorized" {
		contractStruct.AUTHORIZED = signedTimestamp
	}

	json, _ := json.Marshal(contractStruct)
	err := stub.PutState(args[0], json)
	if err != nil {
		return Error(500, "Can't change status")
	}
	return Success(200, "OK", nil)
}

func (cc *SVS) getAllContracts(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	queryString := fmt.Sprintf("{\"selector\": {\"AUTHORIZED\": {\"$regex\": \"^\"}}}")
	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return Error(http.StatusInternalServerError, err.Error())
	}

	defer resultsIterator.Close()
	var results SEARCH
	for resultsIterator.HasNext() {
		it, _ := resultsIterator.Next()
		contract, _ := stub.GetState(it.Key)
		var contractStruct CONTRACT
		json.Unmarshal(contract, &contractStruct)
		results.Contracts = append(results.Contracts, contractStruct)

	}
	resultJson, _ := json.Marshal(results.Contracts)
	return Success(200, "OK", resultJson)

}

func (cc *SVS) createShipment(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	var shipment Shipment
	xmlErr := xml.Unmarshal([]byte(args[0]), &shipment)
	if xmlErr != nil {
		return Error(500, "can't unmarshal shipment xml")
	}
	if value, err := stub.GetState(shipment.ShipmentID); err != nil || value != nil {
		return Error(401, "shipment already exists")
	}	

	json, _ := json.Marshal(shipment)
	if err := stub.PutState(shipment.ShipmentID, json); err != nil {
		return Error(400, "can't create shipment")
	}
	return Success(http.StatusOK, "OK", nil)
}

func (cc *SVS) getAllShipments(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	queryString := fmt.Sprintf("{\"selector\": {\"ShipmentDate\": {\"$regex\": \"^\"}}}")
	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return Error(http.StatusInternalServerError, err.Error())
	}
	defer resultsIterator.Close()

	var results Search_Shipment
	for resultsIterator.HasNext() {
		it, _ := resultsIterator.Next()
		strShipment, _ := stub.GetState(it.Key)		
        var shipmentStruct Shipment
        json.Unmarshal(strShipment, &shipmentStruct)
        results.Shipments = append(results.Shipments, shipmentStruct)

	}
    resultJson, _ := json.Marshal(results.Shipments)
	return Success(200, "OK", resultJson)
}



func (cc *SVS) createOrder(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	var order Order

	xmlErr := xml.Unmarshal([]byte(args[0]), &order)
	if xmlErr != nil {
		return Error(500, "can't unmarshal order xml")
	}

	if value, err := stub.GetState(order.OrderID); err != nil || value != nil {
		return Error(401, "order already exists")
	}

	if order.PaymentTerms == "PK7C" {
		paymentID, _ := newUUID()
		payment := Payment{PaymentID: paymentID, DocumentID: order.OrderID, DocumentType: "order"}
		json, _ := json.Marshal(payment)
		if err := stub.PutState(payment.PaymentID, json); err != nil {
			return Error(400, "can't create prepayment")
		}
	}

	jsonOrder, _ := json.Marshal(order)
	if err := stub.PutState(order.OrderID, jsonOrder); err != nil {
		return Error(400, "can't create order")
	}

	return Success(http.StatusOK, "OK", nil)
}

func (cc *SVS) getAllOrders(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	queryString := fmt.Sprintf("{\"selector\": {\"PaymentTerms\": {\"$regex\": \"^\"}}}")
	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return Error(http.StatusInternalServerError, err.Error())
	}

	defer resultsIterator.Close()

	var resultJson []byte
	for resultsIterator.HasNext() {
		it, _ := resultsIterator.Next()
		strOrder, _ := stub.GetState(it.Key)
		resultJson = append(resultJson, strOrder[:]...)
		// var order Order
		// json.Unmarshal(contract, &contractStruct)
		// results.Contracts = append(results.Contracts, contractStruct)

	}
	// resultJson, _ := json.Marshal(results.Contracts)
	return Success(200, "OK", resultJson)
}

func validateCertData(cert CERTIFICATE, date string) (bool, string) {
	cert.CERT_NUMBER = replaceStr(cert.CERT_NUMBER)
	cert.ORDER = replaceStr(cert.ORDER)
	cert.NAME_SHIP = replaceStr(cert.NAME_SHIP)
	cert.ADDRESS_SHIP = replaceStr(cert.ADDRESS_SHIP)
	cert.PRODUCER = replaceStr(cert.PRODUCER)
	cert.DATE_SHIP = replaceStr(cert.DATE_SHIP)
	cert.WAGON_NUMBER = replaceStr(cert.WAGON_NUMBER)
	cert.CAR_NUMBER = replaceStr(cert.CAR_NUMBER)
	if len(cert.CERT_NUMBER) == 0 {
		return false, "Пустой номер сертификата"
	}
	if len(cert.ORDER) == 0 {
		return false, "Пустой номер приказа на отгрузку"
	}
	if len(cert.NAME_SHIP) == 0 {
		return false, "Пустое наименование получателя"
	}
	if len(cert.ADDRESS_SHIP) == 0 {
		return false, "Пустой адрес получателя"
	}
	if len(cert.PRODUCER) == 0 {
		return false, "Пустое наименование производителя"
	}
	if len(cert.DATE_SHIP) == 0 {
		return false, "Пустая дата отгрузки"
	}
	if len(cert.WAGON_NUMBER) == 0 && len(cert.CAR_NUMBER) == 0 {
		return false, "Пустой номер вагона или номер авто"
	}
	shipDate, _ := time.Parse(dateFormat, cert.DATE_SHIP)
	currentDate, _ := time.Parse(dateFormat, date)
	diff := currentDate.Sub(shipDate)
	if diff.Hours()/24 > dateDiff {
		return false, fmt.Sprintf("Текущая дата превышает дату отгрузки более чем на %.0f дней", dateDiff)
	}
	return true, ""
}

func replaceStr(str string) string {
	value := strings.TrimSpace(str)
	return value
}

func (cc *SVS) uploadXml(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	xmlStr := args[0]
	date := args[1]
	var cert CERTIFICATE
	xmlErr := xml.Unmarshal([]byte(xmlStr), &cert)
	if xmlErr != nil {
		return Error(500, "can't unmarshal xml")
	}
	if err, message := validateCertData(cert, date); !err {
		return Error(400, message)
	}
	cert.STATUS = "Готово к отгрузке"
	for i, _ := range cert.ITEMS {
		cert.ITEMS[i].CERT_NUMBER = cert.CERT_NUMBER
		cert.ITEMS[i].URL = cert.URL
		cert.ITEMS[i].SUBLOT_NUMBER = replaceStr(cert.ITEMS[i].SUBLOT_NUMBER)
		if cert.ITEMS[i].SUBLOT_NUMBER == "" {
			cert.ITEMS[i].STATUS = "Не отгружено"
		} else {
			cert.ITEMS[i].STATUS = "Готово к отгрузке"
		}
	}
	json, jsonErr := json.Marshal(cert)
	if jsonErr != nil {
		return Error(500, "can't marshal json")
	}
	err := stub.PutState(cert.CERT_NUMBER, json)
	if err != nil {
		return Error(500, "can't write data to blockchain")
	}
	return Success(200, "OK", json)
}

func (cc *SVS) status(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	cert, _ := stub.GetState(args[0])
	var certStruct CERTIFICATE
	json.Unmarshal(cert, &certStruct)
	for i, _ := range certStruct.ITEMS {
		if certStruct.ITEMS[i].SUBLOT_NUMBER == args[1] {
			certStruct.ITEMS[i].STATUS = args[2]
		}
	}
	json, _ := json.Marshal(certStruct)
	err := stub.PutState(args[0], json)
	if err != nil {
		return Error(500, "Can't change status")
	}
	return Success(200, "OK", nil)
}

func (cc *SVS) history(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	resultsIterator, err := stub.GetHistoryForKey(args[0])
	if err != nil {
		return Error(http.StatusNotFound, "Not Found")
	}
	defer resultsIterator.Close()

	var results []HISTORY
	for resultsIterator.HasNext() {
		it, _ := resultsIterator.Next()
		var cert CERTIFICATE
		var history HISTORY
		json.Unmarshal(it.Value, &cert)
		history.Certificate = cert
		history.TxId = it.TxId
		history.Timestamp = it.Timestamp.Seconds
		results = append(results, history)
	}
	jsonStr, _ := json.Marshal(results)
	return Success(200, "OK", jsonStr)
}

func (cc *SVS) basicStats(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	queryString := fmt.Sprintf("{\"selector\": {\"CERT_NUMBER\": {\"$regex\": \"^\"}}}")
	resultsIterator, err := stub.GetQueryResult(queryString)
	var stats BASICSTATS
	stats.TotalCert = 0
	stats.TotalLot = 0
	if err != nil {
		return Error(http.StatusInternalServerError, err.Error())
	}
	defer resultsIterator.Close()

	for resultsIterator.HasNext() {
		it, _ := resultsIterator.Next()
		cert, _ := stub.GetState(it.Key)
		var certificate CERTIFICATE
		json.Unmarshal(cert, &certificate)
		stats.TotalCert += 1
		for i, _ := range certificate.ITEMS {
			if certificate.ITEMS[i].SUBLOT_NUMBER != "" {
				stats.TotalLot += 1
			}
		}
	}
	json, _ := json.Marshal(stats)
	return Success(200, "OK", json)
}

func (cc *SVS) stats(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	var stats OMK
	stats.Gost = make(map[string]int)
	stats.Total = 0
	stats.Empty = 0
	stats.Brak = 0
	stats.Way = 0
	stats.Pending = 0
	stats.Accepted = 0
	queryString := fmt.Sprintf("{\"selector\": {\"CERT_NUMBER\": {\"$regex\": \"^\"}}}")
	resultsIterator, err := stub.GetQueryResult(queryString)

	if err != nil {
		return Error(http.StatusInternalServerError, err.Error())
	}
	defer resultsIterator.Close()

	for resultsIterator.HasNext() {
		it, _ := resultsIterator.Next()
		cert, _ := stub.GetState(it.Key)
		var certificate CERTIFICATE
		json.Unmarshal(cert, &certificate)
		for _, item := range certificate.ITEMS {
			stats.Total += 1
			if item.SUBLOT_NUMBER == "" {
				stats.Empty += 1
				stats.EmptyItems = append(stats.EmptyItems, item)
			} else {
				_, ok := stats.Gost[item.GOST]
				if ok {
					stats.Gost[item.GOST] += 1
				} else {
					stats.Gost[item.GOST] = 1
				}
			}
			switch item.STATUS {
			case "Брак":
				stats.Brak += 1
				stats.BrakItems = append(stats.BrakItems, item)
			case "Принято":
				stats.Accepted += 1
			case "В ожидании":
				stats.Pending += 1
			case "В пути":
				stats.Way += 1
			default:
			}
		}
	}
	json, _ := json.Marshal(stats)
	return Success(200, "OK", json)
}

//func (cc *SVS) searchLot(stub shim.ChaincodeStubInterface, args []string) peer.Response {
//    var buffer bytes.Buffer
//    buffer.WriteString("^")
//    buffer.WriteString(args[0])
//    searchStr := buffer.String()
//    buffer.Reset()
//    buffer.WriteString("^")
//	buffer.WriteString(args[1])
//    searchPacket := buffer.String()
//    buffer.Reset()
//    buffer.WriteString("^")
//	buffer.WriteString(args[2])
//    searchStatus := buffer.String()
//    buffer.Reset()
//    buffer.WriteString("^")
//	buffer.WriteString(args[3])
//    searchWagon := buffer.String()
//    buffer.Reset()
//    buffer.WriteString("^")
//	buffer.WriteString(args[4])
//    searchPrikaz := buffer.String()
//    buffer.Reset()
//    buffer.WriteString("^")
//	buffer.WriteString(args[5])
//    searchOrder := buffer.String()
//    buffer.Reset()
//    buffer.WriteString("^")
//	buffer.WriteString(args[6])
//    searchPipe := buffer.String()
//    buffer.Reset()
//
//    queryString := fmt.Sprintf("{\"selector\": {\"CERT_NUMBER\": {\"$regex\": \"%s\"}, \"PACKET_NUMBER\": {\"$regex\": \"%s\"}, \"STATUS\": {\"$regex\": \"%s\"}, \"WAGON_NUMBER\": {\"$regex\": \"%s\"}, \"ORDER\": {\"$regex\": \"%s\"}, \"ORDER_NUMBER\": {\"$regex\": \"%s\"}, \"ITEMS\": {\"$elemMatch\": {\"SUBLOT_NUMBER\": {\"$regex\": \"%s\"}}}}}", searchStr, searchPacket, searchStatus, searchWagon, searchPrikaz, searchOrder, searchPipe)
//
//
//    resultsIterator, err := stub.GetQueryResult(queryString)
//    if err != nil {
//                return Error(http.StatusInternalServerError, err.Error())
//    }
//    defer resultsIterator.Close()
//    var results SEARCH
//    for resultsIterator.HasNext() {
//        it, _ := resultsIterator.Next()
//        cert, _ := stub.GetState(it.Key)
//        var certStruct CERTIFICATE
//        json.Unmarshal(cert, &certStruct)
//        for i, _ := range certStruct.ITEMS {
//        	if certStruct.ITEMS[i].SUBLOT_NUMBER != "" {
//        		results.Items = append(results.Items, certStruct.ITEMS[i])
//        	}
//        }
//    }
//   resultJson, _ := json.Marshal(results.Items)
//    return Success(200, "OK", resultJson)
//}
//
//func (cc *SVS) search(stub shim.ChaincodeStubInterface, args []string) peer.Response {
//    var buffer bytes.Buffer
//    buffer.WriteString("^")
//    buffer.WriteString(args[0])
//    searchStr := buffer.String()
//    buffer.Reset()
//    buffer.WriteString("^")
//	buffer.WriteString(args[1])
//    searchPacket := buffer.String()
//    buffer.Reset()
//    buffer.WriteString("^")
//	buffer.WriteString(args[2])
//    searchStatus := buffer.String()
//    buffer.Reset()
//    buffer.WriteString("^")
//	buffer.WriteString(args[3])
//    searchWagon := buffer.String()
//    buffer.Reset()
//    buffer.WriteString("^")
//	buffer.WriteString(args[4])
//    searchPrikaz := buffer.String()
//    buffer.Reset()
//    buffer.WriteString("^")
//	buffer.WriteString(args[5])
//    searchOrder := buffer.String()
//    buffer.Reset()
//    buffer.WriteString("^")
//	buffer.WriteString(args[6])
//    searchPipe := buffer.String()
//    buffer.Reset()
//
//    queryString := fmt.Sprintf("{\"selector\": {\"CERT_NUMBER\": {\"$regex\": \"%s\"}, \"PACKET_NUMBER\": {\"$regex\": \"%s\"}, \"STATUS\": {\"$regex\": \"%s\"}, \"WAGON_NUMBER\": {\"$regex\": \"%s\"}, \"ORDER\": {\"$regex\": \"%s\"}, \"ORDER_NUMBER\": {\"$regex\": \"%s\"}, \"ITEMS\": {\"$elemMatch\": {\"SUBLOT_NUMBER\": {\"$regex\": \"%s\"}}}}}", searchStr, searchPacket, searchStatus, searchWagon, searchPrikaz, searchOrder, searchPipe)
//
//
//    resultsIterator, err := stub.GetQueryResult(queryString)
//    if err != nil {
//                return Error(http.StatusInternalServerError, err.Error())
//    }
//    defer resultsIterator.Close()
//    var results []CERTIFICATE
//    for resultsIterator.HasNext() {
//        it, _ := resultsIterator.Next()
//        cert, _ := stub.GetState(it.Key)
//        var certStruct CERTIFICATE
//        json.Unmarshal(cert, &certStruct)
//        results = append(results, certStruct)
//    }
//   	resultJson, _ := json.Marshal(results)
//    return Success(200, "OK", resultJson)
//}

func main() {
	if err := shim.Start(new(SVS)); err != nil {
		fmt.Printf("Main: Error starting chaincode: %s", err)
	}
}
