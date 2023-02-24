import "./Navbar.css"

export default function Navbar() {
    return (
        <div className="navbar">
            <div className="logo">
                <h1><a href="#">Kubernetes Sandbox</a></h1>
            </div>

            <ul className="nav">
                <li><a href="#">About</a></li>
                <li><a href="#">Login</a></li>
            </ul>
        </div>
    )
}