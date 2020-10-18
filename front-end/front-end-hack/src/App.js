import React from "react";
import {
  BrowserRouter as Router,
  Switch,
  Route,
  Link
} from "react-router-dom";
import { Home } from "./Home";
import { Login } from "./Login";
import { Account } from "./Account";
import { Register } from "./Register";
import { Edit } from "./Edit";

export default function App() {
  return (
    <Router>
      <div>
        <nav>
          <ul>
            <li>
              <Link to="/">Home</Link>
            </li>
            <li>
              <Link to="/login">Login</Link>
            </li>
            <li>
              <Link to="/account">Account</Link>
            </li>
          </ul>
        </nav>

        <Switch>
          {/* Home */}
          <Route exact path="/">
            <Home />
          </Route>

          {/* Login Page */}
          <Route path="/login">
            <Login />
          </Route>
          {/* Sign up page */}
          <Route path="/register">
            <Register />
          </Route>

          {/* Account page */}
          <Route path="/account">
            <Account />
          </Route>
          {/* Edit page */}
          <Route path="/edit">
            <Edit />
          </Route>

        </Switch>
      </div>
    </Router>
  );
}

