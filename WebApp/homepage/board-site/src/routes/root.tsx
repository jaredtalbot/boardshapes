import boardwalkLogo from "../img/boardwalk.svg";
import boardmeshLogo from "../img/boardmesh.svg";
import githubLogo from "../img/github.svg";
import boardboxLogo from "../img/BoardBox.svg";
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
          <img src={boardmeshLogo} className="logo" alt="Boardmesh logo" />
        </a>
        <a href="/boardbox/" target="_blank">
          <img src={boardboxLogo} className="logo" alt="Boardbox logo" />
        </a>
        <a href="https://github.com/codeJester27/cmps401fa2024" target="_blank">
          <img src={githubLogo} className="logo" alt="Github logo" />
        </a>
      </div>
      <p className="read-the-docs">
        If you wanna know more about the developers, click the link below!
      </p>
      <Link to={`/about`} className="link">
        Click Here!
      </Link>
      <div id="detail">
        <Outlet />
      </div>
    </>
  );
}

export default Root;
