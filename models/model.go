package Models

type User struct {
	Uuid     string    `bson:"uuid,omitempty" json:"uuid,omitempty"`
	Email    string    `bson:"email,omitempty" json:"email,omitempty"`
	Role     string    `bson:"role,omitempty" json:"role,omitempty"`
	Requests []Request ` bson:"requests,omitempty" json:"requests,omitempty"`
}
type VehicleCausing struct {
	Uuid           string `bson:"uuid,omitempty" json:"uuid,omitempty"`
	CarPlateNumber string `bson:"car_plate_number,omitempty" json:"car_plate_number,omitempty"`
	Description    string `bson:"description,omitempty" json:"description,omitempty"`
	Status         string `bson:"status,omitempty" json:"status,omitempty"`
}
type Request struct {
	Uuid           string `bson:"uuid,omitempty" json:"uuid,omitempty"`
	RequestType    string `bson:"request_type,omitempty" json:"request_type,omitempty"`
	CarPlateNumber string `bson:"car_plate_number,omitempty" json:"car_plate_number,omitempty"`
	Description    string `bson:"description,omitempty" json:"description,omitempty"`
	Status         string `bson:"status,omitempty" json:"status,omitempty"`
	Vehicle_wanted bool   `bson:"vehiclewanted,omitempty" json:"vehiclewanted,omitempty"`
}

type GetRequest struct {
	Uuid string `bson:"uuid,omitempty" json:"uuid,omitempty"`
}
