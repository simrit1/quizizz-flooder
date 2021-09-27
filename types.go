package main

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    struct {
		Quizzes []struct {
			CreatedBy struct {
				FirstName string `json:"firstName"`
				LastName  string `json:"lastName"`
				Id        string `json:"_id"`
			} `json:"createdBy"`
		} `json:"quizzes"`
	} `json:"data"`
}

type Config struct {
	Delay      int
	Mode       int
	RoomHash   string
	CustomName string
}

type RoomResp struct {
	Room struct {
		Hash string `json:"hash"`
	} `json:"room"`
}
