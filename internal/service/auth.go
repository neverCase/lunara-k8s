package service

func (s *Service) checkAuth(login *loginRequest) bool {
	groups := map[string]string{
		"admin": "123456",
		"user2": "123456",
	}
	if login.Account == "" || login.Password == "" {
		return false
	}
	if t, ok := groups[login.Account]; ok {
		if t == login.Password {
			return true
		}
	}
	return false
}
