import ApiStatusDisplay from "@/components/api-status-display";

function App() {
  return (
    <div className="min-h-screen bg-gray-950 p-4 flex flex-col gap-4 items-center justify-center">
      <h1 className="text-gray-200 text-2xl font-semibold">
        Server Health Monitor
      </h1>
      <ApiStatusDisplay />
      <footer>
        <p className="text-gray-400 text-sm">
          Made with 💜 by
          <a
            href="https://www.linkedin.com/in/davidcastrosiq/"
            target="_blank"
            rel="noopener noreferrer"
            className="text-blue-500 ml-1"
          >
            David Castro
          </a>
        </p>
      </footer>
    </div>
  );
}

export default App;
