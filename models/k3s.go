package models

type K3sQueryDto struct {
	Kind       string `form:"kind,omitempty"`
	APIVersion string `form:"api_version,omitempty"`

	LabelSelector string `form:"label_selector,omitempty"`
	FieldSelector string `form:"field_selector,omitempty"`

	Watch               bool `form:"watch,omitempty"`
	AllowWatchBookmarks bool `form:"allow_watch_bookmarks,omitempty"`

	ResourceVersion string `form:"resource_version,omitempty"`
	// ResourceVersionMatch ResourceVersionMatch `json:"resourceVersionMatch,omitempty" protobuf:"bytes,10,opt,name=resourceVersionMatch,casttype=ResourceVersionMatch"`

	// TimeoutSeconds *int64 `json:"timeoutSeconds,omitempty" protobuf:"varint,5,opt,name=timeoutSeconds"`

	Limit    int64  `form:"limit,omitempty"`
	Continue string `form:"continue,omitempty"`
}
