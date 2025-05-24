import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import "./index.css";
import { createBrowserRouter, RouterProvider } from "react-router-dom";
import DashBoard from "./component/dashboard/dashboard.tsx";
import Room from "./component/room/room.tsx";
import Lobby from "./component/lobby/lobby.tsx";

const router = createBrowserRouter([
  {
    path: "/",
    element: <DashBoard />,
  },
  {
    path: "/room",
    element: <Room />,
  },
  {
    path: "/lobby/:name",
    element: <Lobby />,
  },
]);
createRoot(document.getElementById("root")!).render(
  // <StrictMode>
    <RouterProvider router={router} />
  // </StrictMode>
);
