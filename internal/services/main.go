package services

import "github.com/samber/do/v2"

var Package = do.Package(
	do.Lazy(NewCaptcha),
	do.Lazy(NewConfig),
	do.Lazy(NewData),
	do.Lazy(NewLogger),
	do.Lazy(NewServer),
)
