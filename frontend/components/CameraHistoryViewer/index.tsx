import React, { useCallback, useEffect, useState } from "react";
import { ImageView } from "./ImageView";
import { useHistoryImages } from "../../hooks/useHistoryImages";
import { Controller } from "./Controller";

interface HistoryResponse {
  files: [HistoryFile];
}

interface HistoryFile {
  path: string;
  timestamp: number;
  invalidated: boolean | undefined;
}

export const CameraHistoryViewer: React.FC = () => {
  const { data } = useHistoryImages();
  const [selected, setSelected] = useState<number | undefined>(undefined);
  const images = data?.files;
  const image = images && selected ? images[selected] : undefined;

  // Control selected image
  const min = 0;
  const max = images ? Math.max(images.length - 1, 0) : 0;

  const selectPrev = useCallback(() => {
    setSelected((selected) =>
      selected !== undefined ? Math.max(min, selected - 1) : undefined
    );
  }, []);

  const selectNext = useCallback(() => {
    setSelected((selected) =>
      selected !== undefined ? Math.min(max, selected + 1) : undefined
    );
  }, [max]);

  const select = useCallback((v: number) => {
    setSelected(v);
  }, []);

  // Set selected on data loaded
  useEffect(() => {
    if (selected === undefined && data) {
      setSelected(data.files.length - 1);
    }
  }, [data, selected]);

  // Manage invalidated
  const [invalidatedTimestamps, setInvalidatedTimestamps] = useState(
    [] as number[]
  );
  const invalidated = image
    ? invalidatedTimestamps.includes(image?.timestamp)
    : false;
  const invalidateImage = useCallback(() => {
    if (image !== undefined) {
      const newList = [...invalidatedTimestamps];
      newList.push(image.timestamp);
      setInvalidatedTimestamps(newList);
    }
  }, [image, invalidatedTimestamps]);

  // Render UI
  return (
    <>
      <ImageView
        src={image?.path}
        invalidated={invalidated}
        invalidate={invalidateImage}
      />
      <Controller
        selectPrev={selectPrev}
        selectNext={selectNext}
        select={select}
        selected={selected}
        min={min}
        max={max}
        disabled={!data}
      />
    </>
  );
};
