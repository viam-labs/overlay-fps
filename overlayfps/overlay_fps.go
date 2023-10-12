package overlayfps

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"math"
	"time"

	"github.com/fogleman/gg"
	"github.com/pkg/errors"

	"github.com/edaniels/golog"
	"go.viam.com/rdk/components/camera"
	"go.viam.com/rdk/resource"
	"go.viam.com/rdk/rimage"
)

// ModelName is the name of the model
const (
	ModelName  = "overlay-fps"
	CountLimit = math.MaxFloat64 / 500.0
)

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
	camera.VideoSource
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
	return camera.FromVideoSource(conf.ResourceName(), o, logger), nil
}

func (o *overlay) Reconfigure(ctx context.Context, deps resource.Dependencies, conf resource.Config) error {
	o.cameraName = ""
	o.VideoSource = nil
	// get the config
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
	vs, ok := cam.(camera.VideoSource)
	if !ok {
		return errors.Wrapf(err, "camera %v is not a video source for %s", cfg.CameraName, ModelName)
	}
	r := &reader{vs, 0.0, 0.0}
	o.VideoSource, err = camera.NewVideoSourceFromReader(ctx, r, nil, props.ImageType)
	if err != nil {
		return err
	}
	return nil
}

// DoCommand simply echos whatever was sent.
func (o *overlay) DoCommand(ctx context.Context, cmd map[string]interface{}) (map[string]interface{}, error) {
	return cmd, nil
}

// Close closes the underlying stream.
func (o *overlay) Close(ctx context.Context) error {
	return o.VideoSource.Close(ctx)
}

type reader struct {
	src      camera.VideoSource
	countDur float64
	avgDur   float64
}

// Read returns the image overlaid  with the FPS
func (r *reader) Read(ctx context.Context) (image.Image, func(), error) {
	if r.countDur >= CountLimit {
		r.countDur = 0.0
	}
	// get image from source camera
	start := time.Now()
	img, release, err := camera.ReadImage(ctx, r.src)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "could not get next source image")
	}
	duration := time.Since(start)
	r.avgDur = (duration.Seconds() + r.avgDur*r.countDur) / (r.countDur + 1.0)
	r.countDur = r.countDur + 1.0
	ovImg := overlayText(img, fmt.Sprintf("avg. FPS: %.2f", 1./r.avgDur))
	return ovImg, release, nil
}

// Close closes the underlying stream.
func (r *reader) Close(ctx context.Context) error {
	return r.src.Close(ctx)
}

// overlayText writes a string in the top of the image.
func overlayText(img image.Image, text string) image.Image {
	gimg := gg.NewContextForImage(img)
	rimage.DrawString(gimg, text, image.Point{30, 30}, color.NRGBA{255, 0, 0, 255}, 30)
	return gimg.Image()
}
