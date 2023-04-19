package config

import "fmt"

func ValidateAllCollaboratorRefsEqual(str1, str2, str3 string) (bool, error) {
	if str1 == str2 && str2 == str3 {
		return true, nil
	}
	return false, fmt.Errorf("collaborator references inside the packages are not equal: %s, %s, %s", str1, str2, str3)
}

func AuthorizeDestinationAccess(collabConfig *CollaborationConfig) (bool, error) {

	return true, nil
}

func AuthorizeTransformationAccess(collabConfig *CollaborationConfig) (bool, error) {

	return true, nil
}
