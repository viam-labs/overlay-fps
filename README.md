# `overlay-fps` modular resource

A virtual camera to overlay the frames-per-second (FPS) of the `get_image()` requests of the underlying camera

![fps-camera-example](https://github.com/viam-labs/overlay-fps/assets/5212232/305a00dd-46fb-41d6-8445-ec7b59d94612)

## Requirements

Before configuring your `overlay-fps` camera, you must [create a machine](https://docs.viam.com/fleet/machines/#add-a-new-machine).

To use the `overlay-fps` camera you also need to [configure a webcam](https://docs.viam.com/components/camera/webcam/) or another [camera](https://docs.viam.com/components/camera/).

## Build and run

To use this module, follow these instructions to [add a module from the Viam Registry](https://docs.viam.com/registry/configure/#add-a-modular-resource-from-the-viam-registry) and select the `viam-labs:camera:overlay-fps` model from the [`overlay-fps` module](https://app.viam.com/module/viam-labs/overlay-fps).

## Configure your `overlay-fps` camera

Navigate to the **Config** tab of your machine's page in [the Viam app](https://app.viam.com/).
Click on the **Components** subtab and click **Create component**.
Select the `camera` type, then select the `viam-labs:camera:overlay-fps` model.
Click **Add module**, then enter a name for your sensor and click **Create**.

On the new component panel, copy and paste the following attribute template into your camera’s **Attributes** box:

```json
{
  "camera_name": "<CAMERA-NAME>"
}
```

Provide your camera's name in the model config.

> [!NOTE]
> For more information, see [Configure a Machine](https://docs.viam.com/manage/configuration/).

### Attributes

The following attributes are available for the `viam-labs:camera:overlay-fps` sensor:

| Name    | Type   | Inclusion    | Description |
| ------- | ------ | ------------ | ----------- |
| `camera_name` | string | **Required** | The name of the camera to add the fps overlay to. |

### Example configuration

```json
{
  "camera_name": "cam1"
}
```

To ensure the source camera starts up before the `overlay-fps` camera, add the source camera in the **Depends on** drop down of the `overlay-fps` camera.

The entire component configuration should resemble this:

```json
{
  "name": "my-overlay-cam",
  "model": "viam-labs:camera:overlay-fps",
  "type": "camera",
  "namespace": "rdk",
  "attributes": {
    "camera_name": "my-webcam"
  },
  "depends_on": [
    "my-webcam"
  ]
}
```

### Next steps

Go to the [**Control** tab](https://docs.viam.com/fleet/machines/#control) and enable the `overlay-fps` camera.
You should now see the camera stream with the FPS overlay.

## Local Development

This module is written in Go.

```bash
go mod tidy
go build
```

## License

Copyright 2021-2023 Viam Inc. <br>
Apache 2.0
