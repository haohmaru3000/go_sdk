package imgprocessing

import "github.com/haohmaru3000/go_sdk/sdkcm"

type Response struct {
	sdkcm.AppError
	Data *sdkcm.Image `json:"data"`
}
