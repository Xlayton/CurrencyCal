import React from "react";
import { Link } from "react-router-dom";

export function Account() {
    return (
        <>
            <h1>Account!</h1>
            <Link to="/edit">Edit Account</Link>
        </>
    );
}
