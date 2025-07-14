import express from 'express';
import path from 'path';
import next from 'next';

const port = parseInt(process.env.PORT || '8080', 10);
const dev = process.env.NODE_ENV !== 'production';
const app = next({ dev, turbo: true });
const handler = app.getRequestHandler();

app.prepare().then(() => {
	const server = express();

	server.use(express.json());
	server.use(express.urlencoded({ extended: true }));

	if (!process.env.TMP_DIR) {
		throw new Error("TMP_DIR env isn't provided");
	}

	server.use(
		'/temp',
		express.static(path.join(__dirname, process.env.TMP_DIR))
	);

	server.all('*path', (req, res) => {
		return handler(req, res);
	});

	server.listen(port, (err) => {
		if (err) throw err;
		console.log(`Listening at port ${port}`);
	});
});
