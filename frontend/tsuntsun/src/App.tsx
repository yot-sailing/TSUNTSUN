import "./App.css";
import { BrowserRouter as Router, Switch, Route } from "react-router-dom";
import Main from "./pages/main";
import Login from "./pages/login";
import AfterLogin from "./pages/afterLogin";

function App() {
  return (
    <div className="App">
      <Router>
        <Switch>
          <Route path="/login">
            <Login />
          </Route>
          <Route path="/after-login">
            <AfterLogin />
          </Route>
          <Route path="/">
            <Main></Main>
          </Route>
        </Switch>
      </Router>
    </div>
  );
}

export default App;
