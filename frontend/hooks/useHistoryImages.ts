import useSWR from "swr";

const fetcher = (...args: Parameters<typeof fetch>) =>
  fetch(...args).then((res) => res.json());

type Response = {
  files: {
    path: string;
    timestamp: number;
  }[];
};

export const useHistoryImages = () => {
  const res = useSWR<Response>("/api/history", fetcher);
  const loading = !res.data && !res.error;
  return { ...res, loading };
};
