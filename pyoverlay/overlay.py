from typing import ClassVar, Mapping, Sequence, Optional, cast, Tuple, List, Any, Dict

from typing_extensions import Self

from viam.module.types import Reconfigurable
from viam.proto.app.robot import ComponentConfig
from viam.proto.common import ResourceName, ResponseMetadata, Geometry
from viam.components.camera import Camera, ViamImage
from viam.resource.types import Model, ModelFamily
from viam.resource.base import ResourceBase
from viam.media.video import CameraMimeType, NamedImage
from viam.media.utils.pil import pil_to_viam_image, viam_to_pil_image
from PIL import ImageDraw, ImageFont

import time


class OverlayCam(Camera, Reconfigurable):
    MODEL: ClassVar[Model] = Model(ModelFamily("viam-labs", "camera"), "overlay-fps")

    def __init__(self, name: str):
        super().__init__(name)

    @classmethod
    def new_cam(
        cls, config: ComponentConfig, dependencies: Mapping[ResourceName, ResourceBase]
    ) -> Self:
        cam = cls(config.name)
        cam.reconfigure(config, dependencies)
        return cam

    @classmethod
    def validate_config(cls, config: ComponentConfig) -> Sequence[str]:
        """Validates JSON configuration"""
        actual_cam = config.attributes.fields["camera_name"].string_value
        if actual_cam == "":
            raise Exception(
                "camera_name attribute is required for a OverlayCam component"
            )
        return [actual_cam]

    def reconfigure(
        self, config: ComponentConfig, dependencies: Mapping[ResourceName, ResourceBase]
    ):
        """Handles attribute reconfiguration"""
        actual_cam_name = config.attributes.fields["camera_name"].string_value
        actual_cam = dependencies[Camera.get_resource_name(actual_cam_name)]
        self.actual_cam = cast(Camera, actual_cam)

    async def get_properties(
        self, *, timeout: Optional[float] = None, **kwargs
    ) -> Camera.Properties:
        """Returns details about the camera"""
        return await self.actual_cam.get_properties()

    async def get_image(
        self,
        mime_type: str = "",
        *,
        extra: Optional[Dict[str, Any]] = None,
        timeout: Optional[float] = None,
        **kwargs
    ) -> ViamImage:
        if mime_type == "":
            mime_type = "image/jpeg"
        start_time = time.time()
        img = await self.actual_cam.get_image(mime_type)
        pil_img = viam_to_pil_image(img)
        end_time = time.time()
        duration = end_time - start_time
        # overlay the fps
        draw = ImageDraw.Draw(pil_img)
        font_size = 30
        # font = ImageFont.load_default()
        font = ImageFont.truetype("FreeMono.ttf", size=font_size)
        position = (30, 30)
        text_color = (255, 0, 0)
        formatted_text = "FPS: {:.2f}".format(1.0 / duration)
        draw.text(position, formatted_text, fill=text_color, font=font)
        return pil_to_viam_image(pil_img, CameraMimeType.from_string(mime_type))

    async def get_images(
        self, *, timeout: Optional[float] = None, **kwargs
    ) -> Tuple[List[NamedImage], ResponseMetadata]:
        raise NotImplementedError

    async def get_point_cloud(
        self,
        *,
        extra: Optional[Dict[str, Any]] = None,
        timeout: Optional[float] = None,
        **kwargs
    ) -> Tuple[bytes, str]:
        raise NotImplementedError

    async def get_geometries(self) -> List[Geometry]:
        raise NotImplementedError
