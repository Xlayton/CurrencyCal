import React from "react";
import { Link } from "react-router-dom";

export function Login() {
    return (
        <>
            <h1>Login!</h1>
            <Link to="/register">Sign Up Here</Link>
        </>
    );
}

