package overlayfps

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"time"

	"github.com/fogleman/gg"
	"github.com/pkg/errors"

	"github.com/edaniels/golog"
	"github.com/viamrobotics/gostream"
	"go.viam.com/rdk/components/camera"
	"go.viam.com/rdk/resource"
	"go.viam.com/rdk/rimage"
)

// ModelName is the name of the model
const ModelName = "overlay-fps"

var (
	// Model is the full resource name of the model
	Model            = resource.NewModel("viam-labs", "camera", ModelName)
	errUnimplemented = errors.New("unimplemented")
)

func init() {
	resource.RegisterComponent(
		camera.API,
		Model,
		resource.Registration[camera.Camera, *Config]{Constructor: newOverlay},
	)
}

// Config specifies which camera and which service should be used to do the overlay
type Config struct {
	CameraName string `json:"camera_name"`
}

// Validate will ensure that the underlying camera is present
func (cfg *Config) Validate(path string) ([]string, error) {
	if cfg.CameraName == "" {
		return nil, fmt.Errorf(`expected "camera_name" attribute for %s %q`, ModelName, path)
	}

	return []string{cfg.CameraName}, nil
}

type overlay struct {
	resource.Named
	stream     gostream.VideoStream
	imgType    camera.ImageType
	cameraName string
	logger     golog.Logger
}

func newOverlay(
	ctx context.Context,
	deps resource.Dependencies,
	conf resource.Config,
	logger golog.Logger,
) (camera.Camera, error) {
	o := &overlay{
		Named:  conf.ResourceName().AsNamed(),
		logger: logger,
	}
	if err := o.Reconfigure(ctx, deps, conf); err != nil {
		return nil, err
	}
	src, err := camera.NewVideoSourceFromReader(ctx, o, nil, o.imgType)
	if err != nil {
		return nil, err
	}
	return camera.FromVideoSource(conf.ResourceName(), src), nil
}

func (o *overlay) Reconfigure(ctx context.Context, deps resource.Dependencies, conf resource.Config) error {
	o.cameraName = ""
	o.stream = nil
	cfg, err := resource.NativeConfig[*Config](conf)
	if err != nil {
		return errors.Errorf("Could not assert proper config for %s", ModelName)
	}
	// get the source camera
	o.cameraName = cfg.CameraName
	cam, err := camera.FromDependencies(deps, cfg.CameraName)
	if err != nil {
		return errors.Wrapf(err, "unable to get camera %v for %s", cfg.CameraName, ModelName)
	}
	props, err := cam.Properties(ctx)
	if err != nil {
		return errors.Wrapf(err, "unable to get camera properties %v for %s", cfg.CameraName, ModelName)
	}
	o.imgType = props.ImageType
	vs, ok := cam.(camera.VideoSource)
	if !ok {
		return errors.Wrapf(err, "camera %v is not a video source for %s", cfg.CameraName, ModelName)
	}
	o.stream = gostream.NewEmbeddedVideoStream(vs)
	return nil
}

// Read returns the image overlaid  with the FPS
func (o *overlay) Read(ctx context.Context) (image.Image, func(), error) {
	// get image from source camera
	start := time.Now()
	img, release, err := o.stream.Next(ctx)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "could not get next source image")
	}
	duration := time.Since(start)
	fps := 1. / duration.Seconds()
	ovImg := overlayText(img, fmt.Sprintf("FPS: %.2f", fps))
	return ovImg, release, nil
}

// overlayText writes a string in the top of the image.
func overlayText(img image.Image, text string) image.Image {
	gimg := gg.NewContextForImage(img)
	rimage.DrawString(gimg, text, image.Point{30, 30}, color.NRGBA{255, 0, 0, 255}, 30)
	return gimg.Image()
}

// DoCommand simply echos whatever was sent.
func (o *overlay) DoCommand(ctx context.Context, cmd map[string]interface{}) (map[string]interface{}, error) {
	return cmd, nil
}

// Close closes the underlying stream.
func (o *overlay) Close(ctx context.Context) error {
	return o.stream.Close(ctx)
}
