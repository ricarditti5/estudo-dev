function MyProps(props) {
    const stringProps = JSON.stringify(props);
    return (
        <div>
            <h1>CHECK OUT MY PROPS OBJECT</h1>
            <h2>{stringProps}</h2>
        </div>
    );
}

export default MyProps;