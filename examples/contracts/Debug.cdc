access(all) contract Debug {


	access(all) struct Foo{
		access(all) let bar: String

		init(bar: String) {
			self.bar=bar
		}
	}

	access(all) event Log(msg: String)

	access(all) fun log(_ msg: String) : String {
		emit Log(msg: msg)
		return msg
	}

}
