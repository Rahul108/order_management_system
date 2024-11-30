package helpers

func CreateAuthorizationFailureMessage() map[string]interface{} {
	return map[string]interface{}{
		"message": "Unauthorized",
		"type":    "error",
		"code":    401,
	}
}
