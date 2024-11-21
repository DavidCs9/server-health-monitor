import { useState, useEffect } from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";

interface ApiData {
  is_up: boolean;
  latency: number;
  server_url: string;
  timestamp: string;
}

export default function ApiStatusDisplay() {
  const [apiData, setApiData] = useState<Array<ApiData> | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchData = async () => {
    try {
      const response = await fetch("http://localhost:8080/health"); // Replace with your actual API endpoint
      if (!response.ok) {
        throw new Error("Failed to fetch data");
      }

      const data = await response.json();
      data.sort((a: ApiData, b: ApiData) => {
        return a.server_url > b.server_url ? 1 : -1;
      });
      setApiData(data);
      setIsLoading(false);
      setError(null);
    } catch (err) {
      console.error(err);
      setError("Error fetching data");
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchData(); // Initial fetch
    const intervalId = setInterval(fetchData, 10000); // Poll every 10 seconds

    return () => clearInterval(intervalId); // Cleanup on unmount
  }, []);

  const formatLatency = (latency: number): string => {
    if (latency < 1000000000) {
      return `${(latency / 1000000).toFixed(2)} ms`;
    } else {
      return `${(latency / 1000000000).toFixed(2)} s`;
    }
  };

  if (isLoading) {
    return <div>Loading...</div>;
  }

  if (error) {
    return <div>Error: {error}</div>;
  }

  return (
    <div className="grid grid-cols-2 gap-4 place-items-center">
      {apiData &&
        apiData.map((data, index) => (
          <Card key={index} className="w-full max-w-md">
            <CardHeader>
              <CardTitle>
                <strong>{data.server_url}</strong>
              </CardTitle>
            </CardHeader>
            <CardContent>
              {data && (
                <div className="space-y-4">
                  <div className="flex items-center justify-between">
                    <span className="text-sm font-medium">Status:</span>
                    <Badge variant={data.is_up ? "success" : "destructive"}>
                      {data.is_up ? "Running" : "Down"}
                    </Badge>
                  </div>
                  <div className="space-y-2">
                    <span className="text-sm font-medium">Latency:</span>
                    <div className="text-3xl font-bold">
                      {formatLatency(data.latency)}
                    </div>
                  </div>
                  <div className="text-sm text-zinc-500 dark:text-zinc-400">
                    <p>
                      Last updated:{" "}
                      <strong>
                        {new Date(data.timestamp).toLocaleString()}
                      </strong>
                    </p>
                  </div>
                </div>
              )}
            </CardContent>
          </Card>
        ))}
    </div>
  );
}
