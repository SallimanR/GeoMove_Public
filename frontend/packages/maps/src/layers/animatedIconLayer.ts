import { IconLayer } from "deck.gl";
import type { MovingPosition } from "../composables/useMovingIconLayer";

export interface AnimatedIconConfig {
  url: string;
  width: number;
  height: number;
}

export function createAnimatedIconLayer(
  data: MovingPosition[],
  icon: AnimatedIconConfig,
  options?: {
    layerId?: string;
    size?: number;
    onClick?: (id: number) => void;
    onHover?: (id: number) => void;
  },
): IconLayer<MovingPosition> {
  const size = options?.size ?? 40;
  const layerId = options?.layerId ?? "moving-icons";

  return new IconLayer<MovingPosition>({
    id: layerId,
    pickable: true,
    data,
    getPosition: (d) => d.position,
    getAngle: (d) => d.bearing,
    getIcon: () => ({ url: icon.url, width: icon.width, height: icon.height }),
    getSize: size,
    billboard: false,
    updateTriggers: {
      getPosition: data,
      getAngle: data,
    },
    onClick: (info) => {
      if (info.object) options?.onClick?.(info.object.id);
    },
    onHover: (info) => {
      if (info.object) options?.onHover?.(info.object.id);
    },
  });
}
