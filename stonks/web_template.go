package stonks

var indexTemplate string = `
<!doctype html>
<html>
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <link href="https://unpkg.com/tailwindcss@^2/dist/tailwind.min.css" rel="stylesheet">
  <title>Crypto Moonshot Gem Finder</title>
</head>
<body class="font-mono align-middle bg-indigo-900 relative overflow-hidden h-screen bg-cover">
	<img src="/images/bg-unsplash.jpg" class="absolute h-full w-full object-cover"/>

	<div class="container mx-auto px-6 md:px-12 relative z-10 flex items-center py-32 xl:py-40">
		<div class="w-full flex flex-col items-center relative z-10">
			<h1 class="font-extrabold text-7xl text-center sm:text-8xl text-white leading-tight mt-4">
				Crypto Moonshot Gems
			</h1>            
		</div>
	</div>

	<div id="gems" class="sm:flex flex-wrap justify-center items-center text-center gap-8">

		{{range .Gems}}
		<!--Card-->
		<div class="z-40 min-h-1/2 max-h-1/2 sm:h-full md:h-full md:h-1/3 w-full sm:w-1/2 md:w-1/2 lg:w-1/5 px-16 py-16 bg-white mt-6 shadow-lg rounded-lg bg-indigo-100 dark:bg-gray-800">
			<div class="flex-shrink-0">
				<div class="flex items-center mx-auto justify-center h-32 w-32 rounded-md">
					<img src="{{.ImageURL}}" class="min-w-full"></img>
				</div>
			</div>
			<h2 class="text-2xl sm:text-xl text-black font-semibold dark:text-white py-4">{{.Name}} (${{.Symbol | ToUpper}})</h2>
			
			<div id="stats" class="sm:flex flex-wrap justify-center items-center text-center gap-2">
				<!-- Price -->
				<div class="h-32 w-1/4 shadow-lg rounded-2xl p-2 bg-indigo-50 dark:bg-purple-200">
					<div class="flex items-center text-left">
						<p class="text-md text-black dark:text-white ml-2">Price</p>
					</div>
					<div class="flex flex-col justify-start">
						<p class="text-gray-700 dark:text-gray-200 text-xl text-center font-bold my-4 slashed-zero">
							{{.CurrentPrice}}
							<span class="text-sm">{{GetCurrency | ToUpper}}</span>
						</p>
					</div>
				</div>
				<!-- MC Rank -->
				<div class="h-32 w-1/4 shadow-lg rounded-2xl p-2 bg-indigo-50 dark:bg-purple-200">
					<div class="flex items-center text-left">
						<p class="text-md text-black dark:text-white ml-2">Market Cap Rank</p>
					</div>
					<div class="flex flex-col justify-start">
						<p class="text-gray-700 dark:text-gray-200 text-xl text-center font-bold my-4 slashed-zero">
							<span class="text-sm">#</span>
							{{.MarketCapRank}}
						</p>
					</div>
				</div>
				<!-- Dev Score -->
				<div class="h-32 w-1/4 shadow-lg rounded-2xl p-2 bg-indigo-50 dark:bg-purple-200">
					<div class="flex items-center text-left">
						<p class="text-md text-black dark:text-white ml-2">Developer Score</p>
					</div>
					<div class="flex flex-col justify-start">
						<p class="text-gray-700 dark:text-gray-200 text-xl text-center font-bold my-4 slashed-zero">
							{{.DeveloperScore}}
							<span class="text-sm">%</span>
						</p>
					</div>
				</div>
			</div> <!-- /stats -->

			<p class="h-64 text-md text-gray-500 dark:text-gray-500 py-4">
				{{.Description}}
			</p>
		</div>
		{{end}}
	</div>

	<div class="absolute bottom-0 right-0 h-16 sm:w-full md:w-1/6 lg:w-1/8 text-green-200">{{.Timestamp}}</div>

</body>
</html>
`
