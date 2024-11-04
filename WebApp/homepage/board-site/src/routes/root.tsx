import boardwalkLogo from "./assetsReal/image2vector2.svg";
import boardmeshLogo from "./assetsReal/image2vector.svg";
import "./root.css";
import { Outlet, Link } from "react-router-dom";

function Root() {
  return (
    <>
      <img src={boardmeshLogo} className="logoMain" alt="Boardwalk logo" />
      <h1>
        Boardmesh: The web API that takes image files and translates them into
        polygonal meshes!
      </h1>
      <div>
        <a href="/boardwalk/" target="_blank">
          <img src={boardwalkLogo} className="logo" alt="Boardwalk logo" />
        </a>
        <a href="/manual/" target="_blank">
          <img
            src={boardmeshLogo}
            className="logo react"
            alt="Boardmesh logo"
          />
        </a>
      </div>
      <p className="read-the-docs">
        If you wanna know more about the developers, click the link below!
      </p>
      <Link to={`/about`}>Click Here!</Link>
      <div id="detail">
        <Outlet />
      </div>
    </>
  );
}

export default Root;
