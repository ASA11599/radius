import { createSignal, onMount } from "solid-js";
import "./App.css";

import { BE_BASE_URL } from "./backend";

function distance(p1, p2) {
	// Haversine formula
	const r = 6371
	const lat2 = p2.latitude * Math.PI / 180
	const lat1 = p1.latitude * Math.PI / 180
	const long2 = p2.longitude * Math.PI / 180
	const long1 = p1.longitude * Math.PI / 180
	const dLat = (lat2 - lat1) / 2
	const dLong = (long2 - long1) / 2
	return 2 * r * Math.asin(Math.sqrt(Math.pow(Math.sin(dLat), 2) + (Math.cos(lat1) * Math.cos(lat2) * Math.pow(Math.sin(dLong), 2))))
}

async function sendPost(lat, lon, msg) {
    try {
        const response = await fetch(`${BE_BASE_URL}/posts`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({
                duration: 30,
                content: msg,
                location: {
                    latitude: lat,
                    longitude: lon
                }
            })
        });
        const posts = await response.json();
        return posts;
    } catch (error) {
        console.error(error);
    }
}

async function getPosts(r, lat, long) {
    try {
        const response = await fetch(`${BE_BASE_URL}/posts?radius=${r}&lat=${lat}&long=${long}`);
        const posts = await response.json();
        return posts;
    } catch (error) {
        console.error(error);
    }
}

function App() {

    const [message, setMessage] = createSignal("");
    const [posts, setPosts] = createSignal([]);

    const [location, setLocation] = createSignal(null);

    const postMessage = async (msg) => {
        await sendPost(location().latitude, location().longitude, msg);
    };

    const refreshPosts = async () => {
        const posts = await getPosts(5, location().latitude, location().longitude);
        setPosts(posts);
    };

    onMount(async () => {
        if (window.isSecureContext) {
            navigator.geolocation.getCurrentPosition(async (position) => {
                setLocation({
                    latitude: position.coords.latitude,
                    longitude: position.coords.longitude
                });
                await refreshPosts();
            }, (err) => console.error(err));
        } else {
            console.info("Insecure context, using position (0, 0)");
            setLocation({
                latitude: 0,
                longitude: 0
            });
            await refreshPosts();
        }
    });

    return (
        <>
        <h1>Radius</h1>
        <div class="card">
            <input onChange={(e) => {setMessage(e.target.value)}} type="text" name="postInput" id="post-input" />
            <br />
            <br />
            <button onClick={async () => {
                if (message() !== "") {
                    try {
                        await postMessage(message());
                        await refreshPosts();
                    } catch (error) {
                        console.error(error);
                    }
                }
            }}>
            Post
            </button>
            <br />
            <br />
            <For each={posts()}>
                {(post, i) => {
                    return (
                        <li>
                            {Math.round(distance(location(), post["location"]) * 100) / 100} km away: {post["content"]}
                        </li>
                    );
                }}
            </For>
        </div>
        </>
    );

}

export default App;
