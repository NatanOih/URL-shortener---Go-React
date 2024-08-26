import { ChangeEvent, FormEvent, useEffect, useState } from "react";
import "./App.css";

interface UrlData {
  url: string;
  clicks: number;
}

type ApiResponse = Record<string, UrlData>;

type ApiPost = {
  newUrl: string;
  customShort: string;
};

function App() {
  const [data, setData] = useState<ApiResponse | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [loadError, setLoadError] = useState<string | null>(null);
  const [formError, setFormError] = useState<string | null>("");

  const [shortenUrlForm, setShortenUrlForm] = useState<ApiPost>({
    newUrl: "",
    customShort: "",
  });

  const BASE_URL = "";

  const fetchData = async () => {
    try {
      const response = await fetch("/api/v1/url-list", {});
      console.log("response", response);

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const result: ApiResponse = await response.json();
      setData(result);
    } catch (err: unknown) {
      if (err instanceof Error && err.name !== "AbortError") {
        setLoadError(err.message);
      }
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchData();
  }, []);

  const formSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    if (!shortenUrlForm.newUrl) {
      setFormError("URL is required");
      return;
    }

    setFormError(null);

    try {
      const response = await fetch(BASE_URL + "/api/v1", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          url: shortenUrlForm.newUrl,
          CustomShort: shortenUrlForm.customShort,
        }),
      });

      if (!response.ok) {
        throw new Error("Failed to shorten URL");
      }

      // Handle successful response
      const result = await response.json();
      console.log("Shortened URL:", result);
    } catch (error) {
      setFormError("Failed to submit the form" + error);
    }

    fetchData();
  };

  if (loading) return <div>Loading...</div>;
  if (loadError) return <div className="text-4xl">Error: {loadError}</div>;

  return (
    <div className="flex flex-col gap-4">
      <h1>URL List</h1>
      {data && (
        <ul className="flex flex-col justify-center gap-4 items-center">
          {Object.entries(data).map(([shortUrl, { url, clicks }]) => (
            <li
              className="flex flex-row gap-10 justify-between w-full items-center text-center"
              key={shortUrl}
            >
              <strong>
                <a target="_blank" href={BASE_URL + "/" + shortUrl}>
                  {shortUrl}
                </a>{" "}
              </strong>
              <span className="max-w-40 truncate">{url}</span>
              <span className=""> Clicks: {clicks}</span>
            </li>
          ))}
        </ul>
      )}

      <form
        className="flex flex-col gap-2 justify-center  p-2 rounded-sm bg-red-800/20 items-center"
        onSubmit={formSubmit}
      >
        <div className="flex justify-between items-center gap-2 w-full">
          <label htmlFor="url" className="flex-1 text-center pr-2">
            URL:
          </label>
          <input
            id="url"
            className="p-1 rounded-sm flex-1"
            type="url"
            value={shortenUrlForm.newUrl}
            placeholder="URL to shorten"
            onChange={(e: ChangeEvent<HTMLInputElement>) => {
              setShortenUrlForm((prev) => {
                return { ...prev, newUrl: e.target.value };
              });
            }}
          />
        </div>
        <div className="flex justify-between items-center gap-2 w-full">
          <label htmlFor="Custom-Short" className="flex-1 text-center pr-2">
            Custom-Short:
          </label>
          <input
            id="Custom-Short"
            className="p-1 rounded-sm flex-1"
            placeholder="Custom short optional"
            value={shortenUrlForm.customShort}
            onChange={(e: ChangeEvent<HTMLInputElement>) => {
              setShortenUrlForm((prev) => {
                return { ...prev, customShort: e.target.value };
              });
            }}
          />
        </div>
        {formError && <p style={{ color: "red" }}>{formError}</p>}
        <button type="submit">Submit</button>
      </form>
    </div>
  );
}

export default App;
