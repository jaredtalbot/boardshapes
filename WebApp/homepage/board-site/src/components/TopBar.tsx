import boardwalkLogo from "../img/boardwalk.svg";
import boardmeshLogo from "../img/boardmesh.svg";
import githubLogo from "../img/github.svg";
import boardboxLogo from "../img/BoardBox.svg";
import Headroom from "react-headroom";

export default function TopBar() {
  return (
    <Headroom className="topbar">
      <nav>
        <div className="link-container">
          <a href="/boardwalk/" target="_blank" title="Boardwalk">
            <img src={boardwalkLogo} className="logo" alt="Boardwalk logo" />
          </a>
          <a href="/manual/" target="_blank" title="User Manual">
            <img src={boardmeshLogo} className="logo" alt="Boardmesh logo" />
          </a>
          <a href="/boardbox/" target="_blank" title="Boardbox">
            <img src={boardboxLogo} className="logo" alt="Boardbox logo" />
          </a>
          <a
            href="https://github.com/codeJester27/cmps401fa2024"
            target="_blank"
            title="GitHub"
            rel="noopener"
          >
            <img src={githubLogo} className="logo" alt="Github logo" />
          </a>
        </div>
      </nav>
    </Headroom>
  );
}
