import { Link } from "react-router-dom";
import "./root.css";
import { useMemo } from "react";

function About() {
  const todayHash = useMemo(() => {
    const today = new Date();
    today.setHours(0, 0, 0, 0);
    return today
      .toISOString()
      .split("")
      .reduce((a, b) => {
        a = (a << 5) - a + b.charCodeAt(0);
        return a & a;
      }, 0);
  }, []);

  return (
    <>
      <div className="body1">
        <h1>About The Devs:</h1>
        <a className="AboutTitle" href={`/boardwalk?unlock=${todayHash}`}>
          Jared
        </a>
        <p className="AboutDescription">
          Jared is a Fifth-Semester Junior majoring in Computer Science. He is
          the project leader!
        </p>
        <p className="AboutTitle">Cohen</p>
        <p className="AboutDescription">
          Cohen is a Fifth-Semester Junior majoring in Computer Science. He made
          this website!
        </p>
        <p className="AboutTitle">Luke</p>
        <p className="AboutDescription">
          Luke is a Fifth-Semester Junior majoring in Computer Science. He's
          hard at work on the CLI tool!
        </p>
        <p className="AboutTitle">Zach</p>
        <p className="AboutDescription">
          Zach is a Fifth-Semester Junior majoring in Computer Science. He made
          the Boardwalk mobile port!
        </p>
        <Link to={`/`} className="link">
          Return Home!
        </Link>
      </div>
    </>
  );
}

export default About;
