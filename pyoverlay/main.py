import asyncio

from viam.components.camera import Camera
from viam.module.module import Module
from viam.resource.registry import Registry, ResourceCreatorRegistration

import overlay

async def main():
    """This function creates and starts a new module, after adding all desired resource models.
    Resource creators must be registered to the resource registry before the module adds the resource model.
    """
    Registry.register_resource_creator(Camera.SUBTYPE, overlay.OverlayCam.MODEL, ResourceCreatorRegistration(overlay.OverlayCam.new_cam, overlay.OverlayCam.validate_config))
    module = Module.from_args()

    module.add_model_from_registry(Camera.SUBTYPE, overlay.OverlayCam.MODEL)
    await module.start()

if __name__ == "__main__":
    asyncio.run(main())
