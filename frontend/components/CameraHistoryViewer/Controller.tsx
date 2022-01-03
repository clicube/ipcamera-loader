import { Box, Button, Slider } from "@mui/material";

export const Controller = ({
  selectPrev,
  selectPrevSkip,
  selectNext,
  selectNextSkip,
  select,
  min,
  max,
  selected = max,
  disabled,
}: {
  selectPrev: () => void;
  selectPrevSkip: () => void;
  selectNext: () => void;
  selectNextSkip: () => void;
  select: (v: number) => void;
  min: number;
  max: number;
  selected?: number;
  disabled: boolean;
}) => {
  return (
    <>
      <Box display="flex" justifyContent="center">
        <Button
          size="large"
          variant="outlined"
          sx={{ m: 1 }}
          onClick={selectPrevSkip}
          disabled={disabled}
        >
          &lt;&lt;
        </Button>
        <Button
          size="large"
          variant="outlined"
          sx={{ m: 1 }}
          onClick={selectPrev}
          disabled={disabled}
        >
          &lt;
        </Button>
        <Button
          size="large"
          variant="outlined"
          sx={{ m: 1 }}
          onClick={selectNext}
          disabled={disabled}
        >
          &gt;
        </Button>
        <Button
          size="large"
          variant="outlined"
          sx={{ m: 1 }}
          onClick={selectNextSkip}
          disabled={disabled}
        >
          &gt;&gt;
        </Button>
      </Box>
      <Slider
        sx={{ my: 2 }}
        disabled={disabled}
        step={1}
        min={min}
        max={max}
        value={selected ?? min}
        onChange={(e, v) => select(v as number)}
      />
    </>
  );
};
