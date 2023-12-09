package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
)

type User struct {
	Name  string
	Email string
}
type JsonUser struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (u *User) Print(level int) {
	ident := strings.Repeat("-", level)
	log.Println(ident, "Username:", u.Name, u.Email)
}
func (u *User) Id() string {
	return fmt.Sprintf("%p", u)
}
func (u *User) MarshalJSON() ([]byte, error) {
	return json.Marshal(&JsonUser{
		ID:    u.Id(),
		Name:  u.Name,
		Email: u.Email,
	})
}
func (u *User) UnmarshalJSON(data []byte) error {
	aux := &JsonUser{}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	u.Name = aux.Name
	u.Email = aux.Email
	load_helper[aux.ID] = u
	log.Println("Added user with id ", aux.ID, u.Name)
	return nil
}

type Record struct {
	Type     string // MX / A / CNAME / TXT / REDIR / SVR
	Name     string // @ / www
	Host     string // IP / address
	Priority int    // Used for MX
	Port     int    // Used for SVR
}
type JsonRecord struct {
	ID       string
	Type     string
	Name     string
	Host     string
	Priority int
	Port     int
}

func (r *Record) Print(level int) {
	ident := strings.Repeat("-", level)
	log.Println(ident, "", r.Type, r.Name, r.Host)
}
func (r *Record) Id() string {
	return fmt.Sprintf("%p", r)
}
func (r *Record) MarshalJSON() ([]byte, error) {
	return json.Marshal(&JsonRecord{
		ID:       r.Id(),
		Name:     r.Name,
		Type:     r.Type,
		Host:     r.Host,
		Priority: r.Priority,
		Port:     r.Port,
	})
}
func (r *Record) UnmarshalJSON(data []byte) error {
	aux := &JsonRecord{}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	r.Name = aux.Name
	r.Type = aux.Type
	r.Host = aux.Host
	r.Priority = aux.Priority
	r.Port = aux.Port
	load_helper[aux.ID] = r
	log.Println("Added record with id ", aux.ID, r.Name)
	return nil
}

type Domain struct {
	Name    string
	User    *User     // User ID
	Records []*Record // Record ID's
}
type JsonDomain struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	User    string   `json:"user"`
	Records []string `json:"records"`
}

func (d *Domain) Print(level int) {
	ident := strings.Repeat("-", level)
	log.Println(ident, "Domain:", d.Name)
	d.User.Print(level + 1)
	log.Println(ident, " Records:")
	for _, r := range d.Records {
		r.Print(level + 2)
	}
}
func (d *Domain) Id() string {
	return fmt.Sprintf("%p", d)
}
func (d *Domain) MarshalJSON() ([]byte, error) {
	var record_ids []string
	for _, r := range d.Records {
		record_ids = append(record_ids, r.Id())
	}
	return json.Marshal(JsonDomain{
		ID:      d.Id(),
		Name:    d.Name,
		User:    d.User.Id(),
		Records: record_ids,
	})
}
func (d *Domain) UnmarshalJSON(data []byte) error {
	log.Println("UnmarshalJSON domain")
	aux := &JsonDomain{}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	d.Name = aux.Name
	d.User = load_helper[aux.User].(*User) // restore pointer to domains user
	for _, record_id := range aux.Records {
		d.Records = append(d.Records, load_helper[record_id].(*Record))
	}
	return nil
}

type State struct {
	Users   map[string]*User
	Records map[string]*Record
	Domains map[string]*Domain
}

func NewState() *State {
	s := &State{}
	s.Users = make(map[string]*User)
	s.Domains = make(map[string]*Domain)
	s.Records = make(map[string]*Record)
	return s
}
func (s *State) Print() {
	log.Println("State:")
	log.Println("Users:")
	for _, u := range s.Users {
		u.Print(1)
	}
	log.Println("Domains:")
	for _, d := range s.Domains {
		d.Print(1)
	}
}
func (s *State) NewUser(name string, email string) *User {
	u := &User{Name: name, Email: email}
	id := fmt.Sprintf("%p", u)
	s.Users[id] = u
	return u
}
func (s *State) NewDomain(user *User, name string) *Domain {
	d := &Domain{Name: name, User: user}
	s.Domains[d.Id()] = d
	return d
}
func (s *State) NewMxRecord(d *Domain, rtype string, name string, host string, priority int) *Record {
	r := &Record{Type: rtype, Name: name, Host: host, Priority: priority}
	d.Records = append(d.Records, r)
	s.Records[r.Id()] = r
	return r
}
func (s *State) FindDomain(name string) (*Domain, error) {
	for _, v := range s.Domains {
		if v.Name == name {
			return v, nil
		}
	}
	return nil, errors.New("Not found")
}
func Save(s *State) (string, error) {
	b, err := json.MarshalIndent(s, "", "    ")
	if err == nil {
		return string(b), nil
	} else {
		log.Println(err)
		return "", err
	}
}

var load_helper map[string]interface{}

func Load(s *State, blob string) {
	load_helper = make(map[string]interface{})
	if err := json.Unmarshal([]byte(blob), s); err != nil {
		log.Println(err)
	} else {
		log.Println("OK")
	}
}

func test_state() {

	s := NewState()
	u := s.NewUser("Ownername", "some@email.com")
	d := s.NewDomain(u, "somedomain.com")
	s.NewMxRecord(d, "MX", "@", "192.168.1.1", 10)
	s.NewMxRecord(d, "A", "www", "192.168.1.1", 0)

	s.Print()

	x, _ := Save(s) // Saved to json string

	log.Println("State saved, the json string is:")
	log.Println(x)

	s2 := NewState() // Create a new empty State
	Load(s2, x)
	s2.Print()

	d, err := s2.FindDomain("somedomain.com")
	if err == nil {
		d.User.Name = "Changed"
	} else {
		log.Println("Error:", err)
	}
	s2.Print()
}

func main() {
	test_state()
}
