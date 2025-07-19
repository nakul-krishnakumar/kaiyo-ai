export default function Home() {
  return (
    <>
      <button
        onClick={() => {
          fetch('http://localhost:8081/auth/login/?provider=google', {
            method: 'GET',
            headers: {
              'Content-Type': 'application/json',
            },
          })
        }}
      >
        Login with Google
      </button>
    </>
  );
}
