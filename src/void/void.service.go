package void

import "grape/src/common/service"

type VoidService struct {
}

func NewVoidService(s *service.CommonService) *VoidService {
	return &VoidService{
		// ProjectRepository: pr.NewProjectRepository(s.DB),
	}
}

func (c *VoidService) Get(src string) (interface{}, error) {
	return nil, nil
}

func (c *VoidService) Save(file interface{}, dst string) error {
	// 	postBody, _ := json.Marshal(map[string]string{
	//       "name":  "Toby",
	//       "email": "Toby@example.com",
	//    })
	//    responseBody := bytes.NewBuffer(postBody)
	// //Leverage Go's HTTP Post function to make request
	//    resp, err := http.Post("https://postman-echo.com/post", "application/json", responseBody)
	// //Handle Error
	//    if err != nil {
	//       log.Fatalf("An Error Occured %v", err)
	//    }
	//    defer resp.Body.Close()
	// //Read the response body

	return nil

}
