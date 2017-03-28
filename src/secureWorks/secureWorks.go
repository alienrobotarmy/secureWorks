package secureWorks

import "fmt"
import "net/http"
import "strings"
import "crypto/tls"
import "encoding/xml"
import "time"
import "bufio"
import "errors"
import "strconv"
import "io/ioutil"

type SOAPFaultEnvelope struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`
	Body    *SOAPFaultBody
}
type SOAPFaultBody struct {
	XMLName xml.Name   `xml:"http://schemas.xmlsoap.org/soap/envelope/ Body"`
	Fault   *SOAPFault `xml:"Fault"`
}
type SOAPFault struct {
	XMLName     xml.Name        `xml:"Fault"`
	FaultCode   string          `xml:"faultcode"`
	FaultString string          `xml:"faultstring"`
	Detail      SOAPFaultDetail `xml:"detail"`
}
type SOAPFaultDetail struct {
	XMLName   xml.Name            `xml:"detail"`
	FaultInfo SOAPFaultDetailInfo `xml:"faultInfo"`
}
type SOAPFaultDetailInfo struct {
	XMLName   xml.Name `xml:"faultInfo"`
	FaultCode string   `xml:"faultCode"`
	Reason    string   `xml:"reason"`
}
type AttachmentResponseEnvelope struct {
	RawXML   string
	Content  string `xml:"Body>getAttachmentResponse>attachment>content"`
	Filename string `xml:"Body>getAttachmentResponse>attachment>filename"`
	Md5Sum   string `xml:"Body>getAttachmentResponse>attachment>md5Sum"`
}
type QueueTicketIdsResponseEnvelope struct {
	RawXML    string
	TicketIds []string `xml:"Body>getQueueTicketIdsResponse>ticketId"`
}
type QueueCountResponseEnvelope struct {
	RawXML string
	Count  int `xml:"Body>getQueueCountResponse>count"`
}
type ContactListResponseEnvelope struct {
	RawXML   string
	Contacts []IdName `xml:"Body>getContactsResponse>getContactList"`
}
type CustomerListResponseEnvelope struct {
	RawXML     string
	ClientInfo []IdName `xml:"Body>getCustomerListResponse>clientInfo"`
}
type DeviceListResponseEnvelope struct {
	RawXML  string
	Devices []DeviceList `xml:"Body>getDeviceListResponse>device"`
}
type DeviceList struct {
	Client      IdName `xml:"client"`
	DeviceAlias string `xml:"deviceAlias"`
	DeviceId    int    `xml:"deviceId"`
	DeviceIp    string `xml:"deviceIp"`
	DeviceName  string `xml:"deviceName"`
	Location    IdName `xml:"location"`
}
type TicketDetailResponseEnvelope struct {
	RawXML string
	Detail Ticket `xml:"Body>getTicketDetailResponse>ticketDetail"`
}
type UpdatesResponseEnvelope struct {
	RawXML  string
	Tickets []Ticket `xml:"Body>getUpdatesResponse>ticket"`
}
type Ticket struct {
	AttachmentId        int       `xml:"attachmentInfo>id"`
	AttachmentName      string    `xml:"attachmentInfo>name"`
	Client              IdName    `xml:"client"`
	Contact             IdName    `xml:"contact"`
	DateClosed          int64     `xml:"dateClosed"`
	DateCreated         int64     `xml:"dateCreated"`
	DateModified        int64     `xml:"dateModified"`
	DetailedDescription string    `xml:"detailedDescription"`
	Devices             IdName    `xml:"devices"`
	EventSource         string    `xml:"eventSource"`
	IsGlobaChild        bool      `xml:"isGlobalChild"`
	IsGlobaParent       bool      `xml:"isGlobalParent"`
	Location            IdName    `xml:"location"`
	Reason              string    `xml:"reason"`
	ResponsibleParty    string    `xml:"responsibleParty"`
	Service             string    `xml:"service"`
	Severity            string    `xml:"severity"`
	Status              string    `xml:"status"`
	SymptomDescription  string    `xml:"symptomDescription"`
	TicketId            string    `xml:"ticketId"`
	TicketType          string    `xml:"ticketType"`
	TicketVersion       string    `xml:"ticketVersion"`
	WorkLogs            []WorkLog `xml:"worklogs"`
}
type WorkLog struct {
	DateCreated int64  `xml:"dateCreated"`
	Description string `xml:"description"`
	Type        string `xml:"type"`
}
type IdName struct {
	Id   int    `xml:"id"`
	Name string `xml:"name"`
}
type Query struct {
	xml.Name   `xml:"Config"`
	UserName   string `xml:"UserName"`
	Password   string `xml:"Password"`
	ClientId   string `xml:"ClientId"`
	LocationId string `xml:"LocationId"`
	ApiUri     string `xml:"ApiUri"`
}

func (s SOAPFaultEnvelope) printErr(function string) {
	fmt.Printf("Server Error (%s)\n\tFaultCode: %s\n\tFaultString: %s\n",
		function, s.Body.Fault.FaultCode, s.Body.Fault.FaultString)
	fmt.Printf("\tFaultInfo\n\t\tFaultCode: %s\n\t\tReason: %s\n",
		s.Body.Fault.Detail.FaultInfo.FaultCode,
		s.Body.Fault.Detail.FaultInfo.Reason)
}
func (s Ticket) PrintDetails() {
	fmt.Printf("Attachment: %s (%d)\n", s.AttachmentName, s.AttachmentId)
	fmt.Printf("Client: %s (%d)\n", s.Client.Name, s.Client.Id)
	fmt.Printf("Contact: %s (%d)\n", s.Contact.Name, s.Contact.Id)
	fmt.Printf("DateClosed: %d\n", s.DateClosed)
	fmt.Printf("DateCreated: %d\n", s.DateCreated)
	fmt.Printf("DateModified: %d\n", s.DateModified)
	fmt.Printf("DetailedDescription: %s\n", s.DetailedDescription)
	fmt.Printf("Devices: %s (%d)\n", s.Devices.Name, s.Devices.Id)
	fmt.Printf("EventSource: %s\n", s.EventSource)
	fmt.Printf("IsGlobaChild: %t\n", s.IsGlobaChild)
	fmt.Printf("IsGlobaParent: %t\n", s.IsGlobaParent)
	fmt.Printf("Location: %s (%d)\n", s.Location.Name, s.Location.Id)
	fmt.Printf("Reason: %s\n", s.Reason)
	fmt.Printf("ResponsibleParty: %s\n", s.ResponsibleParty)
	fmt.Printf("Service: %s\n", s.Service)
	fmt.Printf("Severity: %s\n", s.Severity)
	fmt.Printf("Status: %s\n", s.Status)
	fmt.Printf("SymptomDescription: %s\n", s.SymptomDescription)
	fmt.Printf("TicketId: %s\n", s.TicketId)
	fmt.Printf("TicketType: %s\n", s.TicketType)
	fmt.Printf("TicketVersion: %s\n", s.TicketVersion)
	fmt.Printf("WorkLogs:\n")
}
func (s Ticket) PrintCsv() {
	fmt.Printf("AttachmentName,AttachmentId,ClientName,ClientId,ContactName,ContactId,DateClosed,DateCreated,"+
		"DateModified,DetailedDescription,DeviceName,DeviceId,EventSource,"+
		"IsGlobaChild,IsGlobaParent,LocationName,LocationId,,Reason,ResponsibleParty,"+
		"Service,Severity,"+
		"Status,SymptomDescription,TicketId,TicketType,TicketVersion,WorkLogs\n"+
		"%s,%d,%s,%d,%s,%d,%d,%d,%d,%s,%s,%d,%s,%t,%t,%s,%d,%s,%s,%s,%s,%s,%s,%s,%s,%s\n",
		s.AttachmentName, s.AttachmentId,
		s.Client.Name, s.Client.Id, s.Contact.Name, s.Contact.Id,
		s.DateClosed, s.DateCreated, s.DateModified, s.DetailedDescription,
		s.Devices.Name, s.Devices.Id, s.EventSource, s.IsGlobaChild, s.IsGlobaParent,
		s.Location.Name, s.Location.Id, s.Reason, s.ResponsibleParty, s.Service,
		s.Severity, s.Status, s.SymptomDescription, s.TicketId, s.TicketType, s.TicketVersion)
}
func (s Ticket) PrintWorkLogs() {
	fmt.Printf("DateCreated,Description\n")
	for _, v := range s.WorkLogs {
		fmt.Printf("%d,%s\n", v.DateCreated, v.Description)
	}
}
func GetContactList(q Query) (*ContactListResponseEnvelope, error) {
	SOAPxml := `
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:ser="http://service.ticket.api.mod.secureworks.com/">
   <soapenv:Header/>
   <soapenv:Body>
      <ser:getContacts>
        <userName>` + q.UserName + `</userName>
         <password>` + q.Password + `</password>
         <clientId>` + q.ClientId + `</clientId>
         <locationId>` + q.LocationId + `</locationId>
      </ser:getContacts>
   </soapenv:Body>
</soapenv:Envelope>
`

	x := new(ContactListResponseEnvelope)
	buf, err := makeSOAPrequest(q, SOAPxml, &x)
	x.RawXML = buf
	return x, err
}
func GetCustomerList(q Query) (*CustomerListResponseEnvelope, error) {
	SOAPxml := `
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:ser="http://service.ticket.api.mod.secureworks.com/">
   <soapenv:Header/>
   <soapenv:Body>
      <ser:getCustomerList>
        <userName>` + q.UserName + `</userName>
        <password>` + q.Password + `</password>
      </ser:getCustomerList>
   </soapenv:Body>
</soapenv:Envelope>
`
	x := new(CustomerListResponseEnvelope)
	buf, err := makeSOAPrequest(q, SOAPxml, &x)
	x.RawXML = buf
	return x, err
}
func GetAttachment(q Query, ticketId string, attachmentId string) (*AttachmentResponseEnvelope, error) {
	SOAPxml := `
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:ser="http://service.ticket.api.mod.secureworks.com/">
   <soapenv:Header/>
   <soapenv:Body>
      <ser:getAttachment>
	 <userName>` + q.UserName + `</userName>
         <password>` + q.Password + `</password>
         <ticketId>` + ticketId + `</ticketId>
         <attachmentId>` + attachmentId + `</attachmentId>
      </ser:getAttachment>
   </soapenv:Body>
</soapenv:Envelope>
`
	x := new(AttachmentResponseEnvelope)
	buf, err := makeSOAPrequest(q, SOAPxml, &x)
	x.RawXML = buf
	return x, err
}
func GetTicketDetail(q Query, ticketId string) (*TicketDetailResponseEnvelope, error) {
	SOAPxml := `
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:ser="http://service.ticket.api.mod.secureworks.com/">
   <soapenv:Header/>
   <soapenv:Body>
      <ser:getTicketDetail>
         <userName>` + q.UserName + `</userName>
         <password>` + q.Password + `</password>
         <ticketId>` + ticketId + `</ticketId>
      </ser:getTicketDetail>
   </soapenv:Body>
</soapenv:Envelope>
`
	x := new(TicketDetailResponseEnvelope)
	buf, err := makeSOAPrequest(q, SOAPxml, &x)
	x.RawXML = buf
	return x, err
}
func GetUpdates(q Query, ticketType string, worklogs string, limit int, assignedToCustomer int) (*UpdatesResponseEnvelope, error) {
	SOAPxml := `
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:ser="http://service.ticket.api.mod.secureworks.com/">
   <soapenv:Header/>
   <soapenv:Body>
      <ser:getUpdates>
        <userName>` + q.UserName + `</userName>
         <password>` + q.Password + `</password>
         <ticketType>` + ticketType + `</ticketType>
         <limit>` + strconv.Itoa(limit) + `</limit>
         <worklogs>` + worklogs + `</worklogs>
         <assignedToCustomer>` + assignedToCustomer + `</assignedToCustomer>
      </ser:getUpdates>
   </soapenv:Body>
</soapenv:Envelope>
`
	x := new(UpdatesResponseEnvelope)
	buf, err := makeSOAPrequest(q, SOAPxml, &x)
	x.RawXML = buf
	return x, err
}
func GetQueueTicketIds(q Query, ticketType string, limit int) (*QueueTicketIdsResponseEnvelope, error) {
	SOAPxml := `
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:ser="http://service.ticket.api.mod.secureworks.com/">
   <soapenv:Header/>
   <soapenv:Body>
      <ser:getQueueTicketIds>
         <userName>` + q.UserName + `</userName>
         <password>` + q.Password + `</password>
         <ticketType>` + ticketType + `</ticketType>
         <limit>` + strconv.Itoa(limit) + `</limit>
      </ser:getQueueTicketIds>
   </soapenv:Body>
</soapenv:Envelope>
`
	x := new(QueueTicketIdsResponseEnvelope)
	buf, err := makeSOAPrequest(q, SOAPxml, &x)
	x.RawXML = buf
	return x, err
}
func GetQueueCount(q Query, ticketType string) (*QueueCountResponseEnvelope, error) {
	SOAPxml := `
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:ser="http://service.ticket.api.mod.secureworks.com/">
   <soapenv:Header/>
   <soapenv:Body>
      <ser:getQueueCount>
         <userName>` + q.UserName + `</userName>
         <password>` + q.Password + `</password>
         <ticketType>` + ticketType + `</ticketType>
      </ser:getQueueCount>
   </soapenv:Body>
</soapenv:Envelope>
`
	x := new(QueueCountResponseEnvelope)
	buf, err := makeSOAPrequest(q, SOAPxml, &x)
	x.RawXML = buf
	return x, err
}
func GetDeviceList(q Query) (*DeviceListResponseEnvelope, error) {
	SOAPxml := `
           <soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:ser="http://service.ticket.api.mod.secureworks.com/">
              <soapenv:Header/>
              <soapenv:Body>
                 <ser:getDeviceList>
                    <userName>` + q.UserName + `</userName>
                    <password>` + q.Password + `</password>
                    <clientId>` + q.ClientId + `</clientId>
                    <locationId>` + q.LocationId + `</locationId>
                 </ser:getDeviceList>
              </soapenv:Body>
           </soapenv:Envelope>
	   `
	x := new(DeviceListResponseEnvelope)
	buf, err := makeSOAPrequest(q, SOAPxml, &x)
	x.RawXML = buf
	return x, err
}
func makeSOAPrequest(q Query, SOAPxml string, v interface{}) (string, error) {
	var buf string
	/* Set Insecure */
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr, Timeout: time.Second * 300}

	/* Make SOAP Request */
	resp, err := client.Post(q.ApiUri,
		"Content-Type: text/xml;charset=UTF-8", strings.NewReader(SOAPxml))
	if err != nil {
		fmt.Printf("Error: %q\n", err)
		return buf, errors.New("error")
	}

	/* Read body, store contents in "buf" */
	scanner := bufio.NewScanner(resp.Body)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		buf += scanner.Text()
		buf += " "
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("Error: %q\n", err)
		return buf, errors.New("error")
	}
	resp.Body.Close()

	/* Convert SOAP XML response to struct in v{} */
	err = xml.Unmarshal([]byte(buf), v)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return buf, err
	}

	return buf, nil
}
func ReadConfig(fileName string) (Query, error) {
	q := Query{}
	r, err := ioutil.ReadFile(fileName)
	if err != nil {
		return q, err
	}
	err = xml.Unmarshal(r, &q)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return q, nil
	}
	return q, nil
}
