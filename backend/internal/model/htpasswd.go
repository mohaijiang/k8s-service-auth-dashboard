package model

import "time"

// HtpasswdSecretSummary is the list item for htpasswd secrets.
type HtpasswdSecretSummary struct {
	Name      string    `json:"name"`
	Namespace string    `json:"namespace"`
	UserCount int       `json:"userCount"`
	CreatedAt time.Time `json:"createdAt"`
}

// HtpasswdSecretDetail is the full detail of an htpasswd secret.
type HtpasswdSecretDetail struct {
	Name                   string                 `json:"name"`
	Namespace              string                 `json:"namespace"`
	Users                  []HtpasswdUserEntry    `json:"users"`
	UserCount              int                    `json:"userCount"`
	CreatedAt              time.Time              `json:"createdAt"`
	LinkedSecurityPolicies []LinkedSecurityPolicy `json:"linkedSecurityPolicies"`
}

// HtpasswdUserEntry represents a single user in an htpasswd secret.
type HtpasswdUserEntry struct {
	Username string `json:"username"`
}

// LinkedSecurityPolicy represents a SecurityPolicy that references this htpasswd secret.
type LinkedSecurityPolicy struct {
	Name      string           `json:"name"`
	Namespace string           `json:"namespace"`
	TargetRef PolicyTargetRef  `json:"targetRef"`
}

// PolicyTargetRef is the target reference of a SecurityPolicy.
type PolicyTargetRef struct {
	Name string `json:"name"`
	Kind string `json:"kind"`
}

// CreateHtpasswdRequest is the request body for creating an htpasswd secret.
type CreateHtpasswdRequest struct {
	Name  string                    `json:"name" binding:"required"`
	Users []CreateHtpasswdUserEntry `json:"users" binding:"required,min=1"`
}

// CreateHtpasswdUserEntry is a user entry in the create request.
type CreateHtpasswdUserEntry struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=8,max=128"`
}

// AddHtpasswdUserRequest is the request body for adding a user.
type AddHtpasswdUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=8,max=128"`
}
